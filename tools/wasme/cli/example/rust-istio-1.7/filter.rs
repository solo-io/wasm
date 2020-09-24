// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

use log::debug;
use proxy_wasm::traits::*;
use proxy_wasm::types::*;
use std::time::Duration;

#[no_mangle]
pub fn _start() {
    proxy_wasm::set_http_context(|_, _| -> Box<dyn HttpContext> { Box::new(HttpAuth) });
}

struct HttpAuth;

impl HttpAuth {
    fn fail(&mut self) {
      debug!("auth: allowed");
      self.send_http_response(403, vec![], Some(b"not authorized"));
    }
}

// Implement http functions related to this request.
// This is the core of the filter code.
impl HttpContext for HttpAuth {

    // This callback will be invoked when request headers arrive
    fn on_http_request_headers(&mut self, _: usize) -> Action {
        // get all the request headers
        let headers = self.get_http_request_headers();
        // transform them from Vec<(String,String)> to Vec<(&str,&str)>; as dispatch_http_call needs
        // Vec<(&str,&str)>.
        let ref_headers : Vec<(&str,&str)> = headers.iter().map(|(ref k,ref v)|(k.as_str(),v.as_str())).collect();

        // Dispatch a call to the auth-cluster. Here we assume that envoy's config has a cluster
        // named auth-cluster. We send the auth cluster all our headers, so it has context to
        // perform auth decisions.
        let res = self.dispatch_http_call(
            "auth-cluster", // cluster name
            ref_headers, // headers
            None, // no body
            vec![], // no trailers
            Duration::from_secs(1), // one second timeout
        );

        // If dispatch reutrn an error, fail the request.
        match res {
            Err(_) =>{
                self.fail();
            }
            Ok(_)  => {}
        }

        // the dispatch call is asynchronous. This means it returns immediatly, while the request
        // happens in the background. When the response arrives `on_http_call_response` will be 
        // called. In the mean time, we need to pause the request, so it doesn't continue upstream.
        Action::Pause
    }

    fn on_http_response_headers(&mut self, _: usize) -> Action {
        // Add a header on the response.
        self.set_http_response_header("Hello", Some("world"));
        Action::Continue
    }
}

impl Context for HttpAuth {
    fn on_http_call_response(&mut self, _ : u32, header_size: usize, _: usize, _: usize) {
        // We have a response to the http call!

        // if we have no headers, it means the http call failed. Fail the incoming request as well.
        if header_size == 0 {
            self.fail();
            return;
        }

        // Check if the auth server returned "200", if so call `resume_http_request` so request is
        // sent upstream.
        // Otherwise, fail the incoming request.
        match self.get_http_request_header(":status") {
            Some(ref status) if status == "200"  => {
                self.resume_http_request();
            }
            _ => {
                debug!("auth: not authorized");
                self.fail();
            }
        }
    }
}
