use std::os::raw::c_char;

extern "C" {
    pub fn proxy_get_configuration(
        configuration_ptr: *const *mut u8,
        message_size: *mut usize,
    ) -> WasmResult;

    pub fn proxy_log(level: u32, message_data: *const u8, message_size: usize) -> WasmResult;

    // header values
    pub fn proxy_add_header_map_value(
        hm_type: HeaderMapType,
        key_ptr: c_char,
        key_size: usize,
        value_ptr: c_char,
        value_size: *mut usize,
    ) -> WasmResult;

    pub fn proxy_get_header_map_value(
        hm_type: HeaderMapType,
        key_ptr: c_char,
        key_size: usize,
        value_ptr: *const c_char,
        value_size: *mut usize,
    ) -> WasmResult;

    pub fn proxy_set_header_map_pairs(
        hm_type: HeaderMapType,
        ptr: *const c_char,
        size: usize,
    ) -> WasmResult;

    pub fn proxy_get_header_map_pairs(
        hm_type: HeaderMapType,
        ptr: c_char,
        size: usize,
    ) -> WasmResult;

    pub fn proxy_replace_header_map_value(
        hm_type: HeaderMapType,
        key_ptr: c_char,
        size: usize,
        value_ptr: *const c_char,
        value_size: usize,
    ) -> WasmResult;

    pub fn proxy_remove_header_map_value(
        hm_type: HeaderMapType,
        key_ptr: c_char,
        key_size: usize,
    ) -> WasmResult;

    pub fn proxy_get_header_map_size(hm_type: HeaderMapType, size: *mut usize) -> WasmResult;

    //buffer values
    pub fn proxy_get_buffer_bytes(
        buff_type: BufferType,
        start: u32,
        lengthL: u32,
        ptr: *const c_char,
        size: *mut usize,
    ) -> WasmResult;

    pub fn proxy_get_buffer_status(
        buff_type: BufferType,
        length_ptr: *mut usize,
        flags_ptr: *mut u32,
    ) -> WasmResult;
}

#[repr(C)]
pub enum HeaderMapType {
    RequestHeaders = 0,   // During the onLog callback these are immutable
    RequestTrailers = 1,  // During the onLog callback these are immutable
    ResponseHeaders = 2,  // During the onLog callback these are immutable
    ResponseTrailers = 3, // During the onLog callback these are immutable
    GrpcCreateInitialMetadata = 4,
    GrpcReceiveInitialMetadata = 5,  // Immutable
    GrpcReceiveTrailingMetadata = 6, // Immutable
    HttpCallResponseHeaders = 7,     // Immutable
    HttpCallResponseTrailers = 8,    // Immutable
}
#[repr(C)]
pub enum BufferType {
    HttpRequestBody = 0,       // During the onLog callback these are immutable
    HttpResponseBody = 1,      // During the onLog callback these are immutable
    NetworkDownstreamData = 2, // During the onLog callback these are immutable
    NetworkUpstreamData = 3,   // During the onLog callback these are immutable
    HttpCallResponseBody = 4,  // Immutable
    GrpcReceiveBuffer = 5,     // Immutable
}

#[repr(C)]
pub enum WasmResult {
    Ok = 0,
    // The result could not be found, e.g. a provided key did not appear in a table.
    NotFound = 1,
    // An argument was bad, e.g. did not not conform to the required range.
    BadArgument = 2,
    // A protobuf could not be serialized.
    SerializationFailure = 3,
    // A protobuf could not be parsed.
    ParseFailure = 4,
    // A provided expression (e.g. "foo.bar") was illegal or unrecognized.
    BadExpression = 5,
    // A provided memory range was not legal.
    InvalidMemoryAccess = 6,
    // Data was requested from an empty container.
    Empty = 7,
    // The provided CAS did not match that of the stored data.
    CasMismatch = 8,
    // Returned result was unexpected, e.g. of the incorrect size.
    ResultMismatch = 9,
    // Internal failure: trying check logs of the surrounding system.
    InternalFailure = 10,
    // The connection/stream/pipe was broken/closed unexpectedly.
    BrokenConnection = 11,
}
