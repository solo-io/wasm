
use log::{info, debug};
use std::collections::HashMap;
use std::ffi::CString;
use std::os::raw::{c_char};
// use serde::de;

mod filter;

/// Low-level Proxy-WASM APIs for the host functions.
mod host;
// pub mod filter;

/// Logger that integrates with host's logging system.
pub struct Logger;

static LOGGER: Logger = Logger;

impl Logger {
    pub fn init() -> Result<(), log::SetLoggerError> {
        log::set_logger(&LOGGER).map(|()| log::set_max_level(log::LevelFilter::Trace))
    }

    fn proxywasm_loglevel(level: log::Level) -> host::LogLevel {
        match level {
            log::Level::Trace => host::LogLevel::Trace,
            log::Level::Debug => host::LogLevel::Debug,
            log::Level::Info => host::LogLevel::Info,
            log::Level::Warn => host::LogLevel::Warn,
            log::Level::Error => host::LogLevel::Error,
        }
    }
}

impl log::Log for Logger {
    fn enabled(&self, _metadata: &log::Metadata) -> bool { 
        true
    }

    fn log(&self, record: &log::Record) {
        let level = Logger::proxywasm_loglevel(record.level());
        let message = record.args().to_string();
        unsafe {
            host::proxy_log(level, message.as_ptr(), message.len());
        }
    }

    fn flush(&self) {}
}

/// Always hook into host's logging system.
#[no_mangle]
fn _start() {
    Logger::init().unwrap();
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

// macro_rules! root_context_factory {
//     () => {
        
//     };
// }


static mut CONTEXT_FACTORY_MAP: Vec<(String, &dyn ContextFactory)> = {
    Vec::new()
};

static mut ROOT_CONTEXT_FACTORY_MAP: Vec<(String, &dyn RootContextFactory)> = {
    Vec::new()
};

static mut ROOT_CONTEXT_MAP: Vec<(u32, &dyn RootContext)> = {
    Vec::new()
};

static mut CONTEXT_MAP: Vec<(String, &dyn Context)> = {
    Vec::new()
};


fn get_context(key: String) -> Result<&'static dyn Context, ContextManagerError> {
    Err(ContextManagerError::NoContext)
}

fn get_root_contetxt(key: u32) -> Result<&'static dyn RootContext, ContextManagerError> {
    Err(ContextManagerError::NoRootContext)
}

fn get_context_factory(key: String) -> Result<&'static dyn ContextFactory, ContextManagerError>  {
    Err(ContextManagerError::NoContextFactory)
}

fn get_root_contetxt_factory(key: String) -> Result<&'static dyn RootContextFactory, ContextManagerError>  {
    Err(ContextManagerError::NoRootContextFactory)
}

pub enum ContextManagerError {
    NoContext,
    NoRootContext,
    NoRootContextFactory,
    NoContextFactory,
}


pub trait RootContext {
    fn on_start(&self, configuration_size: u32) -> bool;
    fn on_tick(&self);
    fn validate_configuration(&self, configuration_size: u32) -> bool;
    fn on_configure(&self, configuration_size: u32) -> bool;
    fn on_done(&self) -> bool;
    fn on_queue_ready(&self, token: u32);
}

pub trait Context {}


pub trait RootContextFactory {
    fn root_context(&self) -> RootContext;
}

pub trait ContextFactory {
    fn context(&self) -> Context;
}

static mut ROOT_CONTEXT: &dyn RootContext = &BasicRootContext {};

pub fn get_configuration<'a,T : serde::de::DeserializeOwned>(configuration_size: u32) ->  Option<T> {
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
            // let error: serde_json::Error = 
            //     protobuf::ProtobufError::message_not_initialized("configurtion_ptr is null");
            return None
        }
        config =  Box::from_raw(*configuration_ptr);
        debug!("config {:}, size: {:}", config, *message_size);
        read = std::slice::from_raw_parts(config.as_mut(), *message_size);
    }

    match serde_json::from_slice::<T>(&read) {
        Ok(v) => {
            Some(v)
        },
        Err(e) => {
            debug!("error: {}", e);
            None
        },
    }
}

pub struct BasicRootContext {}

impl RootContext for BasicRootContext {
    fn on_start(&self, _: u32) -> bool {
        info!("on_start");
        true
    }

    fn on_tick(&self) {}

    fn validate_configuration(&self, configuration_size: u32) -> bool {
        info!("validate_configuration");
        let proto_config = get_configuration::<filter::Config>(configuration_size);
        match proto_config {
            Some(v) => {
                info!("validate_config: {:?}", v);
                true
            },
            None => {
                info!("validate_config  error");
                false
            },
        }
    }

    fn on_configure(&self, configuration_size: u32) -> bool {
        info!("on_configure");
        let proto_config = get_configuration::<filter::Config>(configuration_size);
        match proto_config {
            Some(v) => {
                info!("on_configure: {:?}", v);
                true
            },
            None => {
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

struct DecoderFilter<D: StreamDecoder> {
    stream_decoder: D,
}

pub trait StreamDecoder {
    fn on_decode_headers(header_map: &HeaderMap, header_only: bool) -> FilterHeadersStatus;
    fn on_decode_metadata(metadata: &Metadata, header_only: bool) -> FilterMetadataStatus;
    fn on_decode_data(buf: &Buffer, end_stream: bool) -> FilterDataStatus;
    fn on_decode_trailers(header_map: HeaderMap, end_stream: bool) -> FilterTrailersStatus;
}

pub trait StreamEncoder {
    fn on_encode_headers(header_map: u32, header_only: bool) -> FilterHeadersStatus;
    fn on_encode_metadata(metadata: &Metadata, header_only: bool) -> FilterMetadataStatus;
    fn on_encode_data(buf: &Buffer, end_stream: bool) -> FilterDataStatus;
    fn on_encode_trailers(header_map: HeaderMap, end_stream: bool) -> FilterTrailersStatus;
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

/// External APIs for envoy to call into
#[no_mangle]
fn proxy_on_start(root_context_id: u32, configuration_size: u32) -> u32 {
    unsafe {
        ROOT_CONTEXT.on_start(configuration_size) as u32
    }
}
#[no_mangle]
fn proxy_validate_configuration(root_context_id: u32, configuration_size: u32) -> u32 {
    unsafe {
        ROOT_CONTEXT.validate_configuration(configuration_size) as u32
    }
}
#[no_mangle]
fn proxy_on_configure(root_context_id: u32, configuration_size: u32) -> u32 {
    unsafe {
        ROOT_CONTEXT.on_configure(configuration_size) as u32
    }
}
#[no_mangle]
fn proxy_on_tick(root_context_id: u32) {
    unsafe {
        ROOT_CONTEXT.on_tick();
    }
}
#[no_mangle]
fn proxy_on_queue_ready(root_context_id: u32, token: u32) {
    unsafe {
        ROOT_CONTEXT.on_queue_ready(token);
    }
}
#[no_mangle]
fn proxy_on_create(context_id: u32, root_context_id: u32) {}
#[no_mangle]
pub fn proxy_on_context_create(context_id: u32, parent_context_id: u32) {}

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

pub fn proxy_on_http_call_response(
    context_id: u32,
    token: u32,
    headers: u32,
    body_size: u32,
    trailers: u32,
) {}

// gRPC
#[no_mangle]
pub fn proxy_on_grpc_create_initial_metadata(context_id: u32, token: u32, headers: u32) {}
#[no_mangle]
pub fn proxy_on_grpc_receive_initial_metadata(context_id: u32, token: u32, headers: u32) {}
#[no_mangle]
pub fn proxy_on_grpc_trailing_metadata(context_id: u32, token: u32, trailers: u32) {}
#[no_mangle]
pub fn proxy_on_grpc_receive(context_id: u32, token: u32, response_size: u32) {}
#[no_mangle]
pub fn proxy_on_grpc_close(context_id: u32, token: u32, status_code: u32) {}

// The stream/vm has completed.
#[no_mangle]
fn proxy_on_done(context_id: u32) -> u32 {1}
// proxy_on_log occurs after proxy_on_done.
#[no_mangle]
fn proxy_on_log(context_id: u32) {}
// The Context in the proxy has been destroyed and no further calls will be coming.
#[no_mangle]
fn proxy_on_delete(context_id: u32)  {}

