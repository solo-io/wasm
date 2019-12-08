
use crate::host;
use log::{info, debug};
use std::ffi::{CString};
use std::os::raw::{c_char};


pub fn get_configuration<T : serde::de::DeserializeOwned>(configuration_size: u32) ->  Result<T, FFIError> {
  let configuration: *mut u8 = std::ptr::null_mut();
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
          return Err(FFIError::Config)
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
          Err(FFIError::Config)
      },
  }
}

pub fn get_properpty(key: &str) -> Result<host::DataExchange, FFIError> {
  let c_to_print =  match CString::new(key) {
      Ok(v) => v,
      Err(_) => return Err(FFIError::Config)
  };
  let mut value_size: Box<usize> = Box::default();
  let value_ptr: *mut c_char = std::ptr::null_mut();
  let value_ptr_ptr: *const *mut c_char = &(value_ptr);
  unsafe {
      let result = host::proxy_get_property(c_to_print.as_ptr(), key.len(),
          value_ptr_ptr, value_size.as_mut() as *mut usize);
      match result {
              host::WasmResult::Ok => {}
              _ => {
                  debug!("result is not ok {}", result as u32);
                  return Err(FFIError::Config)
              }
          };
      if value_ptr_ptr.is_null() {
          return Err(FFIError::Config)
      }

      // debug!("value_suze: {}", value_size);
      // debug!("str_slice: {:?}", CString::from_raw(value_ptr));
      Ok(host::DataExchange{
          value_ptr: value_ptr,
          value_size: *value_size
      })
  }
}

pub enum FFIError {
  Config,
}
