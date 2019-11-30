use log::{info};
/// Low-level Proxy-WASM APIs for the host functions.
pub mod host;
pub mod filter;

/// Logger that integrates with host's logging system.
pub struct Logger;

static LOGGER: Logger = Logger;

impl Logger {
    pub fn init() -> Result<(), log::SetLoggerError> {
        log::set_logger(&LOGGER).map(|()| log::set_max_level(log::LevelFilter::Trace))
    }

    fn proxywasm_loglevel(level: log::Level) -> u32 {
        match level {
            log::Level::Trace => 0,
            log::Level::Debug => 1,
            log::Level::Info => 2,
            log::Level::Warn => 3,
            log::Level::Error => 4,
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

    
struct Filter {
    decoder: &'static DecoderFilter,
    encoder: &'static EncoderFilter
}

impl Filter {
    pub fn new(decoder: &'static DecoderFilter, encoder: &&'static EncoderFilter) -> Self {
        Filter {decoder, encoder}
    }
}

pub struct HeaderMap {}

pub struct Buffer {}

struct DecoderFilter {}

struct EncoderFilter {}


pub trait StreamDecoder {
    fn on_decode_headers(header_map: &HeaderMap, header_only: bool) -> FilterHeadersStatus;
    fn on_decode_data(buf: &Buffer, end_stream: bool) -> FilterDataStatus;
    fn on_decode_trailers(header_map: HeaderMap, end_stream: bool) -> FilterTrailersStatus;
}

pub trait StreamEncoder {
    fn on_encode_headers(header_map: u32, header_only: bool) -> FilterHeadersStatus;
    fn on_encode_data(buf: &Buffer, end_stream: bool) -> FilterDataStatus;
    fn on_encode_trailers(header_map: HeaderMap, end_stream: bool) -> FilterTrailersStatus;
}

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
extern "C" fn proxy_on_start(root_context_id: u32, configuration_size: u32) -> u32 {
    info!("in proxy_on_start");
    0
}
#[no_mangle]
extern "C" fn proxy_validate_configuration(root_context_id: u32, configuration_size: u32) -> u32 {
    // let b = log::RecordBuilder::new();
    // b.level(log::Level::Debug);
    info!("in proxy_validate_config");
    // log::
    // // b.args(String::from(""));
    // // log::
    // // LOGGER.log(b.build());
    0
}
#[no_mangle]
extern "C" fn proxy_on_configure(root_context_id: u32, configuration_size: u32) -> u32 {
    println!("in proxy_on_configure");
    0
}
#[no_mangle]
extern "C" fn proxy_on_tick(root_context_id: u32) {}
// #[no_mangle]
// extern "C" fn proxy_on_queue_ready(root_context_id: u32, token: u32) {}
#[no_mangle]
extern "C" fn proxy_on_create(context_id: u32, root_context_id: u32) {}
#[no_mangle]
extern "C" fn proxy_on_new_connection(context_id: u32) -> FilterStatus {FilterStatus::Continue}
/// stream decoder
#[no_mangle]
extern "C" fn proxy_on_downstream_data(context_id: u32, data_length: u32, end_of_stream: u32) -> FilterStatus {FilterStatus::StopIteration}
// #[no_mangle]
// extern "C" fn proxy_on_queue_ready(root_context_id: u32, token: u32) {}
