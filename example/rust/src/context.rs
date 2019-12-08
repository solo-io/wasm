use phf::phf_map;    
use log::{info, debug};
use std::collections::HashMap;
use crate::ffi;
use std::sync;
use crate::{CONTEXT_MAP, ROOT_CONTEXT_MAP, CONTEXT_FACTORY_MAP, ROOT_CONTEXT_FACTORY_MAP};

fn ensure_root_context(root_context_id: &u32) -> Result<*mut dyn RootContext, ffi::FFIError>  {

  let prop_cstring = unsafe{ ffi::get_properpty("plugin_root_id")?.cstr() };

  let mut str_buf =  match prop_cstring.to_str() {
      Ok(v) => v,
      Err(e) => {
          debug!("into string error: {}", e);
          return Err(ffi::FFIError::Config)
      }
  };
  debug!("plugin_root_id_str: {}", str_buf);

  match CONTEXT_MAP.get(&root_context_id){
    Some(sync_ctx) => {
      return Ok(sync_ctx.as_root())
    },
    None => {},
  };


  debug!("no root context found");


  

  let root_context_factory = match ROOT_CONTEXT_FACTORY_MAP.get(&mut str_buf) {
      Some(v) => v,
      None => return Err(ffi::FFIError::Config)
  };
  
  let root_ctx = root_context_factory(root_context_id, str_buf);

  debug!("created root_ctx");
  
  let context_factory = match CONTEXT_FACTORY_MAP.get(&mut str_buf) {
      Some(v) => v,
      None => return Err(ffi::FFIError::Config)
  }; 

  let ctx = context_factory(root_context_id, root_ctx.clone());

  debug!("created ctx");

  ROOT_CONTEXT_MAP.insert(str_buf, sync::Arc::new(sync::Mutex::new(root_ctx)));
  self.context_map.insert(*root_context_id, sync::Arc::new(sync::Mutex::new(ctx)));

  Ok(root_ctx)
}

pub enum ContextError {
  Config
}


pub trait RootContext: std::marker::Sync {
  fn on_start(&mut self, configuration_size: u32) -> bool;
  fn on_tick(&self);
  fn validate_configuration(&self, configuration_size: u32) -> bool;
  fn on_configure(&mut self, configuration_size: u32) -> bool;
  fn on_done(&self) -> bool;
  fn on_queue_ready(&self, token: u32);
}

pub trait Context: StreamDecoder + StreamEncoder {
  fn as_root(&self) -> *mut dyn RootContext; 
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

pub type ContextFactory = fn(&u32, *mut dyn RootContext) -> *mut dyn Context;

pub type RootContextFactory = fn(&u32, &str) -> *mut dyn RootContext;

pub struct HeaderMap {
  data: HashMap<String, Vec<String>>,}

pub struct Buffer {}

pub struct Metadata {
  data: HashMap<String, String>,
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
