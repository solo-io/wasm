use log::{info, debug};
use std::collections::HashMap;
use std::sync;
use std::os::raw::{c_char};
use std::ffi::{CString, c_void};

// use serde::de;


/// Low-level Proxy-WASM APIs for the host functions.
mod host;
mod logger;
mod filter;

/// Logger that integrates with host's logging system.
pub struct Logger;

/// Always hook into host's logging system.
#[no_mangle]
fn _start() {
    logger::Logger::init().unwrap();

    ContextManager::add_root_context(1, sync::Arc::new(
        sync::Mutex::new(
            BasicRootContext {
                proto_config: None
            }   
        )
    ));
}

/// Allow host to allocate memory.
#[no_mangle]
fn malloc(size: usize) -> *mut u8 {
    let mut vec: Vec<u8> = Vec::with_capacity(size);
    unsafe {
        vec.set_len(size);
    }
    let slice = vec.into_boxed_slice();
    Box::into_raw(slice) as *mut u8
}

/// Allow host to free memory.
#[no_mangle]
fn free(ptr: *mut u8) {
    if !ptr.is_null() {
        unsafe {
            Box::from_raw(ptr);
        }
    }
}



struct ContextManager {
    root_context_map: HashMap<u32, sync::Arc<sync::Mutex<dyn RootContext>>>,
    context_map: HashMap<String, sync::Arc<sync::Mutex<dyn Context>>>,
    context_factory_map: HashMap<String, sync::Arc<sync::Mutex<dyn ContextFactory>>>,
    root_context_factory_map: HashMap<String, sync::Arc<sync::Mutex<dyn RootContextFactory>>>,
}

static mut CONTEXT_MANAGER: Option<sync::Arc<sync::Mutex<Box<ContextManager>>>> = None;

fn get_context_manager() -> sync::Arc<sync::Mutex<Box<ContextManager>>> {
    unsafe  {
        match CONTEXT_MANAGER.clone() {
            Some(cm) => {
                cm
            }
            None => {
                sync::Arc::new(
                    sync::Mutex::new(
                        Box::new(
                            ContextManager{
                                context_factory_map: HashMap::new(),
                                root_context_factory_map: HashMap::new(),
                                root_context_map: HashMap::new(),
                                context_map: HashMap::new(),
                            }
                        )
                    )
                )
            }
        }
    }
}

impl ContextManager {
    fn add_context(key: String, context: sync::Arc<sync::Mutex<dyn Context>>) {
        get_context_manager().lock().unwrap().context_map.insert(key, context);
    }
    fn add_root_context(key: u32, context: sync::Arc<sync::Mutex<dyn RootContext>>) {
        get_context_manager().lock().unwrap().root_context_map.insert(key, context);
    }
    pub fn add_root_context_factory(key: String, context: sync::Arc<sync::Mutex<dyn RootContextFactory>>) {
        get_context_manager().lock().unwrap().root_context_factory_map.insert(key, context);
    }
    pub fn add_context_factory(key: String, context: sync::Arc<sync::Mutex<dyn ContextFactory>>) {
        get_context_manager().lock().unwrap().context_factory_map.insert(key, context);
    }

    fn ensure_root_context() {}

    fn get_context(&mut self, key: String) -> Result<sync::Arc<sync::Mutex<dyn Context>>, ContextManagerError> {
        match self.context_map.get_mut(&key) {
            Some(v) => {
                Ok(v.clone())
            },
            None => {Err(ContextManagerError::NoContext)}
        }
    }
    
    fn get_root_contetxt(&mut self, key: u32) -> Result<sync::Arc<sync::Mutex<dyn RootContext>>, ContextManagerError> {
        match self.root_context_map.get_mut(&key) {
            Some(v) => {
                Ok(v.clone())
            },
            None => {Err(ContextManagerError::NoContext)}
        }
    }
    
    fn get_context_factory(&mut self, key: String) -> Result<sync::Arc<sync::Mutex<dyn ContextFactory>>, ContextManagerError>  {
        match self.context_factory_map.get_mut(&key) {
            Some(v) => {
                Ok(v.clone())
            },
            None => {Err(ContextManagerError::NoContext)}
        }
    }
    
    fn get_root_contetxt_factory(&mut self, key: String) -> Result<sync::Arc<sync::Mutex<dyn RootContextFactory>>, ContextManagerError>  {
        match self.root_context_factory_map.get_mut(&key) {
            Some(v) => {
                Ok(v.clone())
            },
            None => {Err(ContextManagerError::NoContext)}
        }
    }
}



pub enum ContextManagerError {
    NoContext,
    NoRootContext,
    NoRootContextFactory,
    NoContextFactory,
}


pub trait RootContext: std::marker::Sync {
    fn on_start(&mut self, configuration_size: u32) -> bool;
    fn on_tick(&self);
    fn validate_configuration(&self, configuration_size: u32) -> bool;
    fn on_configure(&mut self, configuration_size: u32) -> bool;
    fn on_done(&self) -> bool;
    fn on_queue_ready(&self, token: u32);
}

pub trait Context: StreamDecoder + StreamEncoder {}


pub trait RootContextFactory {
    fn root_context(&self) -> &dyn RootContext;
}

pub trait ContextFactory {
    fn context(&self) -> &dyn Context;
}

// static mut ROOT_CONTEXT: BasicRootContext = BasicRootContext {
//     proto_config: None
// };

pub fn get_configuration<T : serde::de::DeserializeOwned>(configuration_size: u32) ->  Result<T, EnvoyError> {
    let configuration: *mut u8 = malloc(configuration_size as usize);
    let configuration_ptr: *const *mut u8 = &configuration;
    let mut message_size: Box<usize> = Box::default();
    unsafe {
        let result = host::proxy_get_configuration(configuration_ptr, message_size.as_mut());
        if result != host::WasmResult::Ok {  
            debug!("non-ok result: {:}", result as u32);
        }
    }
    let read;
    let mut config: Box<u8>;
    unsafe {
        if configuration_ptr.is_null() {
            return Err(EnvoyError::ConfigurationError)
        }
        config =  Box::from_raw(*configuration_ptr);
        debug!("config {:}, size: {:}", config, *message_size);
        read = std::slice::from_raw_parts(config.as_mut(), *message_size);
    }

    match serde_json::from_slice::<T>(&read) {
        Ok(v) => {
            Ok(v)
        },
        Err(e) => {
            debug!("error: {}", e);
            Err(EnvoyError::ConfigurationError)
        },
    }
}

pub fn get_properpty(key: &str) -> Result<String, EnvoyError> {
    let c_to_print =  match CString::new(key) {
        Ok(v) => v,
        Err(_) => return Err(EnvoyError::NilPropertyError)
    };
    let mut value_size: Box<usize> = Box::default();
    let mut value_ptr: Box<c_char> = Box::default();
    let value_ptr_ptr: *const *mut c_char = &(value_ptr.as_mut() as *mut c_char);
    unsafe {
        match host::proxy_get_property(c_to_print.as_ptr(), key.len()-1,
            value_ptr_ptr, value_size.as_mut() as *mut usize) {
                host::WasmResult::Ok => {}
                _ => return Err(EnvoyError::NilPropertyError)
            };

    }
    Err(EnvoyError::NilPropertyError)
}

pub enum EnvoyError {
    ConfigurationError,
    NilPropertyError,
}

pub struct BasicRootContext {
    proto_config: Option<filter::Config>
}

impl RootContext for BasicRootContext {
    fn on_start(&mut self, _: u32) -> bool {
        info!("on_start");
        true
    }

    fn on_tick(&self) {}

    fn validate_configuration(&self, configuration_size: u32) -> bool {
        info!("validate_configuration");
        let proto_config = get_configuration::<filter::Config>(configuration_size);
        match proto_config {
            Ok(v) => {
                info!("validate_config: {:?}", v);
                true
            },
            Err(_) => {
                info!("validate_config  error");
                false
            },
        }
    }

    fn on_configure(&mut self, configuration_size: u32) -> bool {
        info!("on_configure");
        let proto_config = get_configuration::<filter::Config>(configuration_size);
        match proto_config {
            Ok(v) => {
                info!("on_configure: {:?}", v);
                self.proto_config = Some(v);
                true
            },
            Err(_) => {
                info!("on_configure error");
                false
            },
        }
    }

    fn on_done(&self) -> bool {true}

    fn on_queue_ready(&self, _: u32) {}

}

pub struct HeaderMap {
    data: HashMap<String, Vec<String>>,}

pub struct Buffer {}

pub struct Metadata {
    data: HashMap<String, String>,
}


pub trait StreamDecoder {
    fn on_decode_headers(&self, header_map: &HeaderMap, header_only: bool) -> FilterHeadersStatus;
    fn on_decode_metadata(&self, metadata: &Metadata, header_only: bool) -> FilterMetadataStatus;
    fn on_decode_data(&self, buf: &Buffer, end_stream: bool) -> FilterDataStatus;
    fn on_decode_trailers(&self, header_map: HeaderMap, end_stream: bool) -> FilterTrailersStatus;
}

pub trait StreamEncoder {
    fn on_encode_headers(&self, header_map: u32, header_only: bool) -> FilterHeadersStatus;
    fn on_encode_metadata(&self, metadata: &Metadata, header_only: bool) -> FilterMetadataStatus;
    fn on_encode_data(&self, buf: &Buffer, end_stream: bool) -> FilterDataStatus;
    fn on_encode_trailers(&self, header_map: HeaderMap, end_stream: bool) -> FilterTrailersStatus;
}

#[repr(C)]
pub enum FilterStatus { Continue = 0, StopIteration = 1 }
#[repr(C)]
pub enum FilterHeadersStatus { Continue = 0, StopIteration = 1 }
#[repr(C)]
pub enum FilterMetadataStatus { Continue = 0 }
#[repr(C)]
pub enum FilterTrailersStatus { Continue = 0, StopIteration = 1 }
#[repr(C)]
pub enum FilterDataStatus {
    Continue = 0,
    StopIterationAndBuffer = 1,
    StopIterationAndWatermark = 2,
    StopIterationNoBuffer = 3
}

#[no_mangle]
fn proxy_on_create(context_id: u32, root_context_id: u32) {}

/// External APIs for envoy to call into
#[no_mangle]
fn proxy_on_start(root_context_id: u32, configuration_size: u32) -> u32 {
    match get_context_manager().lock().unwrap().get_root_contetxt(root_context_id) {
        Ok(ctx) => {
            ctx.lock().unwrap().on_start(configuration_size) as u32
        }
        Err(_) => {false as u32}
    }
}
#[no_mangle]
fn proxy_validate_configuration(root_context_id: u32, configuration_size: u32) -> u32 {
    match get_context_manager().lock().unwrap().get_root_contetxt(root_context_id) {
        Ok(ctx) => {
            ctx.lock().unwrap().validate_configuration(configuration_size) as u32
        }
        Err(_) => {false as u32}
    }
}
#[no_mangle]
fn proxy_on_configure(root_context_id: u32, configuration_size: u32) -> u32 {
    match get_context_manager().lock().unwrap().get_root_contetxt(root_context_id) {
        Ok(ctx) => {
            ctx.lock().unwrap().on_configure(configuration_size) as u32
        }
        Err(_) => {false as u32}
    }
}
#[no_mangle]
fn proxy_on_tick(root_context_id: u32) {
    match get_context_manager().lock().unwrap().get_root_contetxt(root_context_id) {
        Ok(ctx) => {
            ctx.lock().unwrap().on_tick();
        }
        Err(_) => {}
    }
}
#[no_mangle]
fn proxy_on_queue_ready(root_context_id: u32, token: u32) {
    match get_context_manager().lock().unwrap().get_root_contetxt(root_context_id) {
        Ok(ctx) => {
            ctx.lock().unwrap().on_queue_ready(token);
        }
        Err(_) => {}
    }
}
#[no_mangle]
fn proxy_on_new_connection(context_id: u32) -> FilterStatus {FilterStatus::Continue}
/// stream decoder
#[no_mangle]
fn proxy_on_downstream_data(context_id: u32, data_length: u32, end_of_stream: u32) -> FilterStatus {FilterStatus::StopIteration}
#[no_mangle]
fn proxy_on_upstream_data(context_id: u32, data_length: u32, end_of_stream: u32) -> FilterStatus {FilterStatus::StopIteration}

#[no_mangle]
fn proxy_on_downstream_connection_close(context_id: u32, peer_type: u32) {}
#[no_mangle]
fn proxy_on_upstream_connection_close(context_id: u32, peer_type: u32) {}

#[no_mangle]
fn proxy_on_request_headers(context_id: u32, headers: u32) -> FilterHeadersStatus {FilterHeadersStatus::Continue}
#[no_mangle]
fn proxy_on_request_metadata(context_id: u32, elements: u32) -> FilterMetadataStatus {FilterMetadataStatus::Continue}
#[no_mangle]
fn proxy_on_request_body(context_id: u32, body_buffer_length: u32, end_of_stream: u32) -> FilterDataStatus {FilterDataStatus::Continue}
#[no_mangle]
fn proxy_on_request_trailers(context_id: u32, trailers: u32) -> FilterTrailersStatus {FilterTrailersStatus::Continue}

#[no_mangle]
fn proxy_on_response_headers(context_id: u32, headers: u32) -> FilterHeadersStatus {FilterHeadersStatus::Continue}
#[no_mangle]
fn proxy_on_response_metadata(context_id: u32, elements: u32) -> FilterMetadataStatus {FilterMetadataStatus::Continue}
#[no_mangle]
fn proxy_on_response_body(context_id: u32, body_buffer_length: u32, end_of_stream: u32) -> FilterDataStatus {FilterDataStatus::Continue}
#[no_mangle]
fn proxy_on_response_trailers(context_id: u32, trailers: u32) -> FilterTrailersStatus {FilterTrailersStatus::Continue}


#[no_mangle]
fn proxy_on_done(context_id: u32) -> u32 {1}
#[no_mangle]
fn proxy_on_log(context_id: u32) {}
#[no_mangle]
fn proxy_on_delete(context_id: u32)  {}


