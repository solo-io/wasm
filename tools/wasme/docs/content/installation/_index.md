---
title: "Installation"
description: "Installing the WebAssembly Hub CLI"
weight: 3
---

To install the WebAssembly Hub CLI (`wasme`), simply run the following:

```bash
curl -sL https://run.solo.io/wasme/install | sh
export PATH=$HOME/.wasme/bin:$PATH
```

Verify that `wasme` installed correctly:
```bash
wasme --version
```

```
wasme version 0.0.16
```

Great! You're all set to start building filters. If you're just getting started with the WebAssembly Hub, check out the [Getting Started Tutorial]({{< versioned_link_path fromRoot="/tutorial_code/getting_started">}})
