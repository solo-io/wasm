use std::os::raw::{c_char, c_void};

extern "C" {
    pub fn proxy_get_configuration(
        configuration_ptr: *const *mut u8,
        message_size: *mut usize,
    ) -> WasmResult;

    pub fn proxy_log(
        level: LogLevel,
        logMessage: *const u8,
        messageSize: usize,
    ) -> WasmResult;

    // Headers
    pub fn proxy_add_header_map_value(
        type_: HeaderMapType,
        key_ptr: *const c_char,
        key_size: usize,
        value_ptr: *const c_char,
        value_size: usize,
    ) -> WasmResult;

    pub fn proxy_get_header_map_value(
        type_: HeaderMapType,
        key_ptr: *const c_char,
        key_size: usize,
        value_ptr: *mut *const c_char,
        value_size: *mut usize,
    ) -> WasmResult;

    pub fn proxy_get_header_map_pairs(
        type_: HeaderMapType,
        ptr: *mut *const c_char,
        size: *mut usize,
    ) -> WasmResult;

    pub fn proxy_set_header_map_pairs(
        type_: HeaderMapType,
        ptr: *const c_char,
        size: usize,
    ) -> WasmResult;

    pub fn proxy_replace_header_map_value(
        type_: HeaderMapType,
        key_ptr: *const c_char,
        key_size: usize,
        value_ptr: *const c_char,
        value_size: usize,
    ) -> WasmResult;

    pub fn proxy_remove_header_map_value(
        type_: HeaderMapType,
        key_ptr: *const c_char,
        key_size: usize,
    ) -> WasmResult;

    pub fn proxy_get_header_map_size(type_: HeaderMapType, size: *mut usize) -> WasmResult;

    //buffer values
    pub fn proxy_get_buffer_bytes(
        type_: BufferType,
        start: u32,
        length: u32,
        ptr: *mut *const c_char,
        size: *mut usize,
    ) -> WasmResult;

    pub fn proxy_get_buffer_status(
        type_: BufferType,
        length_ptr: *mut usize,
        flags_ptr: *mut u32,
    ) -> WasmResult;

    //HTTP
    pub fn proxy_http_call(
        uri_ptr: *const c_char,
        uri_size: usize,
        header_pairs_ptr: *mut c_void,
        header_pairs_size: usize,
        body_ptr: *const c_char,
        body_size: usize,
        trailer_pairs_ptr: *mut c_void,
        trailer_pairs_size: usize,
        timeout_milliseconds: u32,
        token_ptr: *mut u32,
    ) -> WasmResult;

    // gRPC
    pub fn proxy_grpc_call(
        service_ptr: *const c_char,
        service_size: usize,
        service_name_ptr: *const c_char,
        service_name_size: usize,
        method_name_ptr: *const c_char,
        method_name_size: usize,
        request_ptr: *const c_char,
        request_size: usize,
        timeout_milliseconds: u32,
        token_ptr: *mut u32,
    ) -> WasmResult;

    pub fn proxy_grpc_stream(
        service_ptr: *const c_char,
        service_size: usize,
        service_name_ptr: *const c_char,
        service_name_size: usize,
        method_name_ptr: *const c_char,
        method_name_size: usize,
        token_ptr: *mut u32,
    ) -> WasmResult;

    pub fn proxy_grpc_cancel(token: u32) -> WasmResult;

    pub fn proxy_grpc_close(token: u32) -> WasmResult;

    pub fn proxy_grpc_send(
        token: u32,
        message_ptr: *const c_char,
        message_size: usize,
        end_stream: u32,
    ) -> WasmResult;

    // Metrics
    pub fn proxy_define_metric(
        type_: MetricType,
        name_ptr: *const c_char,
        name_size: usize,
        metric_id: *mut u32,
    ) -> WasmResult;

    pub fn proxy_increment_metric(metric_id: u32, offset: i64) -> WasmResult;

    pub fn proxy_record_metric(metric_id: u32, value: u64) -> WasmResult;

    pub fn proxy_get_metric(metric_id: u32, result: *mut u64) -> WasmResult;

    // Results status details for any previous ABI call and onGrpcClose.
    pub fn proxy_get_status(
        status_code_ptr: *mut u32,
        message_ptr: *mut *const c_char,
        message_size: *mut usize,
    ) -> WasmResult;

    // Timer (must be called from a root context, e.g. onStart, onTick).
    pub fn proxy_set_tick_period_milliseconds(millisecond: u32) -> WasmResult;

    // Time
    pub fn proxy_get_current_time_nanoseconds(nanoseconds: *mut u64) -> WasmResult;

    // State accessors
    pub fn proxy_get_property(
        path_ptr: *const c_char,
        path_size: usize,
        value_ptr_ptr: *mut *const c_char,
        value_size_ptr: *mut usize,
    ) -> WasmResult;

    pub fn proxy_set_property(
        path_ptr: *const c_char,
        path_size: usize,
        value_ptr: *const c_char,
        value_size: usize,
    ) -> WasmResult;

    // Continue/Reply/Route
    pub fn proxy_continue_request() -> WasmResult;

    pub fn proxy_continue_response() -> WasmResult;

    pub fn proxy_send_local_response(
        response_code: u32,
        response_code_details_ptr: *const c_char,
        response_code_details_size: usize,
        body_ptr: *const c_char,
        body_size: usize,
        additional_response_header_pairs_ptr: *const c_char,
        additional_response_header_pairs_size: usize,
        grpc_status: u32,
    ) -> WasmResult;

    pub fn proxy_clear_route_cache() -> WasmResult;

    // SharedData
    // Returns: Ok, NotFound
    pub fn proxy_get_shared_data(
        key_ptr: *const c_char,
        key_size: usize,
        value_ptr: *mut *const c_char,
        value_size: *mut usize,
        cas: *mut u32,
    ) -> WasmResult;

    //  If cas != 0 and cas != the current cas for 'key' return false, otherwise set the value and
    //  return true.
    // Returns: Ok, CasMismatch
    pub fn proxy_set_shared_data(
        key_ptr: *const c_char,
        key_size: usize,
        value_ptr: *const c_char,
        value_size: usize,
        cas: u32,
    ) -> WasmResult;

    // SharedQueue
    // Note: Registering the same queue_name will overwrite the old registration while preseving any
    // pending data. Consequently it should typically be followed by a call to
    // proxy_dequeue_shared_queue. Returns: Ok
    pub fn proxy_register_shared_queue(
        queue_name_ptr: *const c_char,
        queue_name_size: usize,
        token: *mut u32,
    ) -> WasmResult;

    // Returns: Ok, NotFound
    pub fn proxy_resolve_shared_queue(
        vm_id: *const c_char,
        vm_id_size: usize,
        queue_name_ptr: *const c_char,
        queue_name_size: usize,
        token: *mut u32,
    ) -> WasmResult;

    // Returns Ok, Empty, NotFound (token not registered).
    pub fn proxy_dequeue_shared_queue(
        token: u32,
        data_ptr: *mut *const c_char,
        data_size: *mut usize,
    ) -> WasmResult;

    // Returns false if the queue was not found and the data was not enqueued.
    pub fn proxy_enqueue_shared_queue(
        token: u32,
        data_ptr: *const c_char,
        data_size: usize,
    ) -> WasmResult;

    // System -- ??? No clue what these actually are for
    pub fn proxy_set_effective_context(effective_context_id: u32) -> WasmResult;
    pub fn proxy_done() -> WasmResult;
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
#[derive(PartialEq)]
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

#[repr(C)]
pub enum LogLevel {
    Trace = 0,
    Debug = 1,
    Info = 2,
    Warn = 3,
    Error = 4,
    Critical = 5
}

#[repr(C)]
pub enum MetricType {
    Counter = 0,
    Gauge = 1,
    Histogram = 2
}