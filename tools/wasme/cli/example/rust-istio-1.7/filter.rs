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

use log::trace;
use proxy_wasm::traits::*;
use proxy_wasm::types::*;
use std::time::Duration;

#[no_mangle]
pub fn _start() {
    proxy_wasm::set_http_context(|_, _| -> Box<dyn HttpContext> { Box::new(HttpAuth) });
}

struct HttpAuth;

impl HttpContext for HttpAuth {
    fn on_http_request_headers(&mut self, _: usize) -> Action {
        self.dispatch_http_call(
            "auth-cluster",
            get_http_request_headers(),
            None,
            vec![],
            Duration::from_secs(1),
        )
        .unwrap();
        Action::Pause
    }
}

impl Context for HttpAuth {
    fn on_http_call_response(&mut self, _: u32, _: usize, _: usize, _: usize) {
        if get_http_request_header(":status") == "200" {
            debug!("auth: allowed");
            self.resume_http_request();
            return;
        }
        debug!("auth: not authorized");
        self.send_http_response(403, vec![], Some(b"not authorized"));
    }
}
