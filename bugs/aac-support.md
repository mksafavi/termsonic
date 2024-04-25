---
title: AAC file format support
reported-on: 2024-04-24
status: open
link: "https://git.sixfoisneuf.fr/list/thread/20240423114537.49426-1-nicholas%40fitzroydale.org.html#2BBA4D8E-A52D-49F3-9B89-8D9610612AFE@fitzroydale.org"
---

The underlying library used by termsonic to handle sound, [gopxl/beep](https://github.com/gopxl/beep), does not support AAC files.

Related issue in their repository: [goplx/beep#108](https://github.com/gopxl/beep/issues/108).

There does not seem to be an agreed upon solution at the moment: there are no pure-Go AAC decoders, and they are reluctant on having to integrate with a library which uses CGO (and so am I).