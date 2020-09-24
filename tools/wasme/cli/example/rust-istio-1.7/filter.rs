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

impl HttpContext for HttpAuth {
    fn on_http_request_headers(&mut self, _: usize) -> Action {
        let headers = self.get_http_request_headers();
        let ref_headers : Vec<(&str,&str)> = headers.iter().map(|(ref k,ref v)|(k.as_str(),v.as_str())).collect();
        let res = self.dispatch_http_call(
            "auth-cluster",
            ref_headers,
            None,
            vec![],
            Duration::from_secs(1),
        );
        match res{
            Err(_) =>{
                self.fail();
            }
            Ok(_)  => {}
        }
        Action::Pause
    }

    fn on_http_response_headers(&mut self, _: usize) -> Action {
        self.set_http_response_header("Hello", Some("world"));
        Action::Continue
    }
}

impl Context for HttpAuth {
    fn on_http_call_response(&mut self, _: u32, _: usize, _: usize, _: usize) {
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
