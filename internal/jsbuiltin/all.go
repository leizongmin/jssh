package jsbuiltin

var modules []JsModule

func init() {
	// builtin_0.js
	modules = append(modules, JsModule{File: "builtin_0.js", Code: "Y29uc3QgZ2xvYmFsID0gZ2xvYmFsVGhpcyB8fCB0aGlzOwoKewogIGNvbnN0IHJlbW92ZVNoZWJhbmdMaW5lID0gKGRhdGEpID0+IHsKICAgIGlmICghZGF0YS5zdGFydHNXaXRoKCIjISIpKSByZXR1cm4gZGF0YTsKICAgIHJldHVybiBkYXRhLnJlcGxhY2UoL14jIVteXG5dKi8sICIiKTsKICB9OwoKICBjb25zdCByZXNvbHZlV2l0aEV4dGVuc2lvbiA9IChuYW1lKSA9PiB7CiAgICBjb25zdCBleHRlbnNpb24gPSBbIi5qc29uIiwgIi5qcyJdOwogICAgaWYgKGZzLmV4aXN0KG5hbWUpKSB7CiAgICAgIGlmIChmcy5zdGF0KG5hbWUpLmlzZGlyKSB7CiAgICAgICAgLy8g5aaC5p6c5piv55uu5b2V77yM5bCd6K+VICR7bmFtZX0vcGFja2FnZS5qc29uCiAgICAgICAgY29uc3QgcGtnRmlsZSA9IHBhdGguam9pbihuYW1lLCAicGFja2FnZS5qc29uIik7CiAgICAgICAgaWYgKGZzLmV4aXN0KHBrZ0ZpbGUpKSB7CiAgICAgICAgICBjb25zdCBwa2cgPSBsb2FkSnNvbk1vZHVsZShwa2dGaWxlKTsKICAgICAgICAgIGlmIChwa2cubWFpbikgewogICAgICAgICAgICByZXR1cm4gcmVzb2x2ZVdpdGhFeHRlbnNpb24ocGF0aC5qb2luKG5hbWUsIHBrZy5tYWluKSk7CiAgICAgICAgICB9CiAgICAgICAgfQogICAgICAgIC8vIOWGjeWwneivlSAke25hbWV9L2luZGV4LmpzLCAke25hbWV9L2luZGV4Lmpzb24KICAgICAgICBjb25zdCBpbmRleEZpbGUgPSBwYXRoLmpvaW4obmFtZSwgImluZGV4Iik7CiAgICAgICAgaWYgKGZzLmV4aXN0KGluZGV4RmlsZSkpIHsKICAgICAgICAgIHJldHVybiBpbmRleEZpbGU7CiAgICAgICAgfQogICAgICAgIGZvciAoY29uc3QgZXh0IG9mIGV4dGVuc2lvbikgewogICAgICAgICAgaWYgKGZzLmV4aXN0KGluZGV4RmlsZSArIGV4dCkpIHsKICAgICAgICAgICAgcmV0dXJuIGluZGV4RmlsZSArIGV4dDsKICAgICAgICAgIH0KICAgICAgICB9CiAgICAgIH0gZWxzZSB7CiAgICAgICAgLy8g5paH5Lu25YiZ55u05o6l6L+U5ZueCiAgICAgICAgcmV0dXJuIG5hbWU7CiAgICAgIH0KICAgIH0gZWxzZSB7CiAgICAgIC8vIOWmguaenOaWh+S7tuS4jeWtmOWcqO+8jOWwneivlSAke25hbWV9LmpzLCAke25hbWV9Lmpzb24KICAgICAgZm9yIChjb25zdCBleHQgb2YgZXh0ZW5zaW9uKSB7CiAgICAgICAgaWYgKGZzLmV4aXN0KG5hbWUgKyBleHQpKSB7CiAgICAgICAgICByZXR1cm4gbmFtZSArIGV4dDsKICAgICAgICB9CiAgICAgIH0KICAgICAgLy8g5YaN5bCd6K+VICR7bmFtZX0vaW5kZXguanMsICR7bmFtZX0vaW5kZXguanNvbgogICAgICBjb25zdCBpbmRleEZpbGUgPSBwYXRoLmpvaW4obmFtZSwgImluZGV4Iik7CiAgICAgIGlmIChmcy5leGlzdChpbmRleEZpbGUpKSB7CiAgICAgICAgcmV0dXJuIGluZGV4RmlsZTsKICAgICAgfQogICAgICBmb3IgKGNvbnN0IGV4dCBvZiBleHRlbnNpb24pIHsKICAgICAgICBpZiAoZnMuZXhpc3QoaW5kZXhGaWxlICsgZXh0KSkgewogICAgICAgICAgcmV0dXJuIGluZGV4RmlsZSArIGV4dDsKICAgICAgICB9CiAgICAgIH0KICAgIH0KICB9OwoKICBjb25zdCByZXNvbHZlTW9kdWxlUGF0aCA9IChuYW1lLCBkaXIpID0+IHsKICAgIGlmIChuYW1lID09PSAiLiIgfHwgbmFtZS5zdGFydHNXaXRoKCIvIikgfHwgbmFtZS5zdGFydHNXaXRoKCIuLyIpKSB7CiAgICAgIHJldHVybiByZXNvbHZlV2l0aEV4dGVuc2lvbihwYXRoLmpvaW4oZGlyLCBuYW1lKSk7CiAgICB9CiAgICBjb25zdCBwYXRocyA9IFtdOwogICAgbGV0IGQgPSBkaXI7CiAgICB3aGlsZSAodHJ1ZSkgewogICAgICBsZXQgcCA9IHBhdGguYWJzKHBhdGguam9pbihkLCAibm9kZV9tb2R1bGVzIikpOwogICAgICBwYXRocy5wdXNoKHApOwogICAgICBjb25zdCBkMiA9IHBhdGguZGlyKGQpOwogICAgICBpZiAoZDIgPT09IGQpIHsKICAgICAgICBicmVhazsKICAgICAgfSBlbHNlIHsKICAgICAgICBkID0gZDI7CiAgICAgIH0KICAgIH0KICAgIGZvciAoY29uc3QgcCBvZiBwYXRocykgewogICAgICBjb25zdCByZXQgPSByZXNvbHZlV2l0aEV4dGVuc2lvbihwYXRoLmpvaW4ocCwgbmFtZSkpOwogICAgICBpZiAocmV0KSByZXR1cm4gcmV0OwogICAgfQogIH07CgogIGNvbnN0IHJlcXVpcmVtb2R1bGUgPSAobmFtZSwgZGlyID0gX19kaXJuYW1lKSA9PiB7CiAgICBpZiAodHlwZW9mIG5hbWUgIT09ICJzdHJpbmciKSB7CiAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoYG1vZHVsZSBuYW1lIGV4cGVjdGVkIHN0cmluZyB0eXBlYCk7CiAgICB9CiAgICBpZiAoIW5hbWUpIHsKICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihgZW1wdHkgbW9kdWxlIG5hbWVgKTsKICAgIH0KICAgIGlmICghZGlyKSB7CiAgICAgIHRocm93IG5ldyBUeXBlRXJyb3IoYGVtcHR5IG1vZHVsZSBkaXJgKTsKICAgIH0KCiAgICBsZXQgZmlsZSA9IHJlc29sdmVNb2R1bGVQYXRoKG5hbWUsIGRpcik7CiAgICBpZiAoIWZpbGUpIHsKICAgICAgdGhyb3cgbmV3IEVycm9yKGBjYW5ub3QgcmVzb2x2ZSBtb2R1bGUgIiR7bmFtZX0iIG9uIHBhdGggIiR7ZGlyfSJgKTsKICAgIH0KICAgIGZpbGUgPSBwYXRoLmFicyhmaWxlKTsKCiAgICB0cnkgewogICAgICBpZiAoZmlsZS5lbmRzV2l0aCgiLmpzb24iKSkgewogICAgICAgIHJldHVybiBsb2FkSnNvbk1vZHVsZShmaWxlKTsKICAgICAgfSBlbHNlIHsKICAgICAgICBjb25zdCBjb250ZW50ID0gZnMucmVhZGZpbGUoZmlsZSk7CiAgICAgICAgcmV0dXJuIGxvYWRKc01vZHVsZShmaWxlLCBwYXRoLmRpcihmaWxlKSwgcmVtb3ZlU2hlYmFuZ0xpbmUoY29udGVudCkpOwogICAgICB9CiAgICB9IGNhdGNoIChlcnIpIHsKICAgICAgY29uc3QgZXJyMiA9IG5ldyBFcnJvcihgY2Fubm90IGxvYWQgbW9kdWxlICIke25hbWV9IjogJHtlcnIubWVzc2FnZX1gKTsKICAgICAgZXJyMi5tb2R1bGVOYW1lID0gbmFtZTsKICAgICAgZXJyMi5yZXNvbHZlZEZpbGVuYW1lID0gZmlsZTsKICAgICAgZXJyMi5vcmlnaW5FcnJvciA9IGVycjsKICAgICAgdGhyb3cgZXJyMjsKICAgIH0KICB9OwoKICBjb25zdCBsb2FkSnNvbk1vZHVsZSA9IChmaWxlbmFtZSkgPT4gewogICAgaWYgKHJlcXVpcmUuY2FjaGVbZmlsZW5hbWVdKSB7CiAgICAgIHJldHVybiByZXF1aXJlLmNhY2hlW2ZpbGVuYW1lXTsKICAgIH0KICAgIHJldHVybiAocmVxdWlyZS5jYWNoZVtmaWxlbmFtZV0gPSBKU09OLnBhcnNlKGZzLnJlYWRmaWxlKGZpbGVuYW1lKSkpOwogIH07CgogIGNvbnN0IGxvYWRKc01vZHVsZSA9IChmaWxlbmFtZSwgZGlybmFtZSwgY29udGVudCkgPT4gewogICAgaWYgKHJlcXVpcmUuY2FjaGVbZmlsZW5hbWVdKSB7CiAgICAgIHJldHVybiByZXF1aXJlLmNhY2hlW2ZpbGVuYW1lXTsKICAgIH0KCiAgICBjb25zdCB3cmFwcGVkID0gYAooZnVuY3Rpb24gKHJlcXVpcmUsIG1vZHVsZSwgX19kaXJuYW1lLCBfX2ZpbGVuYW1lKSB7IHZhciBleHBvcnRzID0gbW9kdWxlLmV4cG9ydHM7ICR7Y29udGVudH0KcmV0dXJuIG1vZHVsZTsKfSkoZnVuY3Rpb24gcmVxdWlyZShuYW1lKSB7CiAgcmV0dXJuIHJlcXVpcmVtb2R1bGUobmFtZSwgIiR7ZGlybmFtZX0iKTsKfSwge2V4cG9ydHM6e30scGFyZW50OnRoaXN9LCAiJHtkaXJuYW1lfSIsICIke2ZpbGVuYW1lfSIpCmAudHJpbUxlZnQoKTsKICAgIHJldHVybiAocmVxdWlyZS5jYWNoZVtfX2ZpbGVuYW1lXSA9IGV2YWxmaWxlKF9fZmlsZW5hbWUsIHdyYXBwZWQpLmV4cG9ydHMpOwogIH07CgogIGNvbnN0IHJlcXVpcmUgPSAobmFtZSkgPT4gewogICAgcmV0dXJuIHJlcXVpcmVtb2R1bGUobmFtZSwgX19kaXJuYW1lKTsKICB9OwoKICByZXF1aXJlLmNhY2hlID0ge307CgogIGdsb2JhbC5yZXF1aXJlID0gcmVxdWlyZTsKICBnbG9iYWwucmVxdWlyZW1vZHVsZSA9IHJlcXVpcmVtb2R1bGU7Cn0K"})

	// builtin_cli.js
	modules = append(modules, JsModule{File: "builtin_cli.js", Code: "Y29uc3QgY2xpID0ge307Cgp7CiAgY29uc3QgX2FyZ3MgPSAoY2xpLl9hcmdzID0gW10pOwogIGNvbnN0IF9vcHRzID0gKGNsaS5fb3B0cyA9IHt9KTsKCiAgZnVuY3Rpb24gZ2V0RmxhZ05hbWUocykgewogICAgaWYgKHMuc3RhcnRzV2l0aCgiLS0iKSkgewogICAgICByZXR1cm4gcy5zbGljZSgyKTsKICAgIH0KICAgIGlmIChzLnN0YXJ0c1dpdGgoIi0iKSkgewogICAgICByZXR1cm4gcy5zbGljZSgxKTsKICAgIH0KICB9CgogIGZvciAobGV0IGkgPSAyOyBpIDwgX19hcmdzLmxlbmd0aDsgaSsrKSB7CiAgICBjb25zdCB2ID0gX19hcmdzW2ldOwogICAgY29uc3QgdjIgPSBfX2FyZ3NbaSArIDFdOwogICAgaWYgKHYuc3RhcnRzV2l0aCgiLSIpKSB7CiAgICAgIGNvbnN0IHIgPSB2Lm1hdGNoKC9eLS0/KFtcd1wtX10rKT0oLiopJC8pOwogICAgICBpZiAocikgewogICAgICAgIF9vcHRzW3JbMV1dID0gclsyXTsKICAgICAgfSBlbHNlIHsKICAgICAgICBpZiAodjIgIT09IHVuZGVmaW5lZCkgewogICAgICAgICAgaWYgKHYyLnN0YXJ0c1dpdGgoIi0iKSkgewogICAgICAgICAgICBfb3B0c1tnZXRGbGFnTmFtZSh2KV0gPSB0cnVlOwogICAgICAgICAgfSBlbHNlIHsKICAgICAgICAgICAgX29wdHNbZ2V0RmxhZ05hbWUodildID0gdjI7CiAgICAgICAgICAgIGkrKzsKICAgICAgICAgIH0KICAgICAgICB9IGVsc2UgewogICAgICAgICAgX29wdHNbZ2V0RmxhZ05hbWUodildID0gdHJ1ZTsKICAgICAgICB9CiAgICAgIH0KICAgIH0gZWxzZSB7CiAgICAgIF9hcmdzLnB1c2godik7CiAgICB9CiAgfQoKICBjbGkuZ2V0ID0gZnVuY3Rpb24gZ2V0KG4pIHsKICAgIGlmICh0eXBlb2YgbiA9PT0gIm51bWJlciIpIHsKICAgICAgcmV0dXJuIF9hcmdzW25dOwogICAgfSBlbHNlIHsKICAgICAgcmV0dXJuIF9vcHRzW25dOwogICAgfQogIH07CgogIGNsaS5ib29sID0gZnVuY3Rpb24gYm9vbChuKSB7CiAgICBpZiAoX29wdHNbbl0gPT09IGZhbHNlIHx8IF9vcHRzW25dID09PSB1bmRlZmluZWQpIHJldHVybiBmYWxzZTsKICAgIGlmIChfb3B0c1tuXSA9PT0gdHJ1ZSkgcmV0dXJuIHRydWU7CiAgICBjb25zdCBzID0gX29wdHNbbl0udG9Mb3dlckNhc2UoKTsKICAgIHJldHVybiAhKHMgPT09ICIwIiB8fCBzID09PSAiZiIgfHwgcyA9PT0gImZhbHNlIik7CiAgfTsKCiAgY2xpLmFyZ3MgPSBmdW5jdGlvbiBhcmdzKCkgewogICAgcmV0dXJuIFsuLi5fYXJnc107CiAgfTsKCiAgY2xpLm9wdHMgPSBmdW5jdGlvbiBvcHRzKCkgewogICAgcmV0dXJuIHsgLi4uX29wdHMgfTsKICB9OwoKICBjbGkucHJvbXB0ID0gZnVuY3Rpb24gcHJvbXB0KG1lc3NhZ2UpIHsKICAgIGlmIChtZXNzYWdlKSBwcmludChtZXNzYWdlKTsKICAgIHJldHVybiByZWFkbGluZSgpOwogIH07CgogIGNsaS5fc3ViY29tbWFuZCA9IHt9OwoKICBjbGkuc3ViY29tbWFuZCA9IGZ1bmN0aW9uIHN1YmNvbW1hbmQobmFtZSwgY2FsbGJhY2spIHsKICAgIGlmICh0eXBlb2YgY2FsbGJhY2sgIT09IGBmdW5jdGlvbmApIHsKICAgICAgdGhyb3cgbmV3IFR5cGVFcnJvcihgY2FsbGJhY2sgZXhwZWN0ZWQgYSBmdW5jdGlvbmApOwogICAgfQogICAgaWYgKGNsaS5fc3ViY29tbWFuZFtuYW1lXSkgewogICAgICB0aHJvdyBuZXcgRXJyb3IoYHN1YmNvbW1hbmQgJHtuYW1lfSBpcyBhbHJlYWR5IHJlZ2lzdGVyZWRgKTsKICAgIH0KICAgIGNsaS5fc3ViY29tbWFuZFtuYW1lXSA9IGNhbGxiYWNrOwogIH07CgogIGNsaS5zdWJjb21tYW5kc3RhcnQgPSBmdW5jdGlvbiBzdWJjb21tYW5kc3RhcnQoKSB7CiAgICBjb25zdCBuYW1lID0gY2xpLmdldCgwKTsKICAgIGlmIChjbGkuX3N1YmNvbW1hbmRbbmFtZV0pIHsKICAgICAgcmV0dXJuIGNsaS5fc3ViY29tbWFuZFtuYW1lXSgpOwogICAgfQogICAgaWYgKGNsaS5fc3ViY29tbWFuZFtgKmBdKSB7CiAgICAgIHJldHVybiBjbGkuX3N1YmNvbW1hbmRbYCpgXSgpOwogICAgfQogICAgdGhyb3cgbmV3IEVycm9yKGB1bnJlY29nbml6ZWQgc3ViY29tbWFuZCAke25hbWV9YCk7CiAgfTsKfQo="})

	// builtin_console.js
	modules = append(modules, JsModule{File: "builtin_console.js", Code: "Y29uc3QgY29uc29sZSA9IHt9OwoKewogIGZ1bmN0aW9uIHByaW50VmFycyguLi52YXJzKSB7CiAgICBwcmludGxuKHZhcnMubWFwKCh2KSA9PiBTdHJpbmcodikpLmpvaW4oIiAiKSk7CiAgfQoKICBjb25zb2xlLmxvZyA9IGZ1bmN0aW9uIGxvZyguLi5hcmdzKSB7CiAgICBwcmludFZhcnMoLi4uYXJncyk7CiAgfTsKCiAgY29uc29sZS5lcnJvciA9IGZ1bmN0aW9uIGVycm9yKC4uLmFyZ3MpIHsKICAgIHByaW50VmFycyguLi5hcmdzKTsKICB9Owp9Cg=="})

	// builtin_exec.js
	modules = append(modules, JsModule{File: "builtin_exec.js", Code: "ZnVuY3Rpb24gZXhlYzEoY21kLCBlbnYgPSB7fSkgewogIHJldHVybiBleGVjKGNtZCwgZW52LCAxKTsKfQoKZnVuY3Rpb24gZXhlYzIoY21kLCBlbnYgPSB7fSkgewogIHJldHVybiBleGVjKGNtZCwgZW52LCAyKTsKfQoKc3NoLmV4ZWMxID0gZnVuY3Rpb24gZXhlYzEoY21kLCBlbnYgPSB7fSkgewogIHJldHVybiBzc2guZXhlYyhjbWQsIGVudiwgMSk7Cn07Cgpzc2guZXhlYzIgPSBmdW5jdGlvbiBleGVjMihjbWQsIGVudiA9IHt9KSB7CiAgcmV0dXJuIHNzaC5leGVjKGNtZCwgZW52LCAyKTsKfTsK"})

	// builtin_global.js
	modules = append(modules, JsModule{File: "builtin_global.js", Code: "ZnVuY3Rpb24gcmFuZG9tc3RyaW5nKHNpemUsIGNoYXJzKSB7CiAgc2l6ZSA9IHNpemUgfHwgNjsKICBjaGFycyA9CiAgICBjaGFycyB8fCAiQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVphYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ejAxMjM0NTY3ODkiOwogIGNvbnN0IG1heCA9IGNoYXJzLmxlbmd0aDsKICBsZXQgcmV0ID0gIiI7CiAgZm9yIChsZXQgaSA9IDA7IGkgPCBzaXplOyBpKyspIHsKICAgIHJldCArPSBjaGFycy5jaGFyQXQoTWF0aC5mbG9vcihNYXRoLnJhbmRvbSgpICogbWF4KSk7CiAgfQogIHJldHVybiByZXQ7Cn0KCmZ1bmN0aW9uIGZvcm1hdGRhdGUoZm9ybWF0LCB0aW1lc3RhbXApIHsKICAvLyAgZGlzY3VzcyBhdDogaHR0cDovL3BocGpzLm9yZy9mdW5jdGlvbnMvZGF0ZS8KICAvLyAgIGV4YW1wbGUgMTogZGF0ZSgnSDptOnMgXFxtIFxcaVxccyBcXG1cXG9cXG5cXHRcXGgnLCAxMDYyNDAyNDAwKTsKICAvLyAgIHJldHVybnMgMTogJzA5OjA5OjQwIG0gaXMgbW9udGgnCiAgLy8gICBleGFtcGxlIDI6IGRhdGUoJ0YgaiwgWSwgZzppIGEnLCAxMDYyNDYyNDAwKTsKICAvLyAgIHJldHVybnMgMjogJ1NlcHRlbWJlciAyLCAyMDAzLCAyOjI2IGFtJwogIC8vICAgZXhhbXBsZSAzOiBkYXRlKCdZIFcgbycsIDEwNjI0NjI0MDApOwogIC8vICAgcmV0dXJucyAzOiAnMjAwMyAzNiAyMDAzJwogIC8vICAgZXhhbXBsZSA0OiB4ID0gZGF0ZSgnWSBtIGQnLCAobmV3IERhdGUoKSkuZ2V0VGltZSgpLzEwMDApOwogIC8vICAgZXhhbXBsZSA0OiAoeCsnJykubGVuZ3RoID09IDEwIC8vIDIwMDkgMDEgMDkKICAvLyAgIHJldHVybnMgNDogdHJ1ZQogIC8vICAgZXhhbXBsZSA1OiBkYXRlKCdXJywgMTEwNDUzNDAwMCk7CiAgLy8gICByZXR1cm5zIDU6ICc1MycKICAvLyAgIGV4YW1wbGUgNjogZGF0ZSgnQiB0JywgMTEwNDUzNDAwMCk7CiAgLy8gICByZXR1cm5zIDY6ICc5OTkgMzEnCiAgLy8gICBleGFtcGxlIDc6IGRhdGUoJ1cgVScsIDEyOTM3NTAwMDAuODIpOyAvLyAyMDEwLTEyLTMxCiAgLy8gICByZXR1cm5zIDc6ICc1MiAxMjkzNzUwMDAwJwogIC8vICAgZXhhbXBsZSA4OiBkYXRlKCdXJywgMTI5MzgzNjQwMCk7IC8vIDIwMTEtMDEtMDEKICAvLyAgIHJldHVybnMgODogJzUyJwogIC8vICAgZXhhbXBsZSA5OiBkYXRlKCdXIFktbS1kJywgMTI5Mzk3NDA1NCk7IC8vIDIwMTEtMDEtMDIKICAvLyAgIHJldHVybnMgOTogJzUyIDIwMTEtMDEtMDInCgogIGxldCBqc2RhdGUsIGY7CiAgLy8gS2VlcCB0aGlzIGhlcmUgKHdvcmtzLCBidXQgZm9yIGNvZGUgY29tbWVudGVkLW91dCBiZWxvdyBmb3IgZmlsZSBzaXplIHJlYXNvbnMpCiAgLy8gdmFyIHRhbD0gW107CiAgY29uc3QgdHh0X3dvcmRzID0gWwogICAgIlN1biIsCiAgICAiTW9uIiwKICAgICJUdWVzIiwKICAgICJXZWRuZXMiLAogICAgIlRodXJzIiwKICAgICJGcmkiLAogICAgIlNhdHVyIiwKICAgICJKYW51YXJ5IiwKICAgICJGZWJydWFyeSIsCiAgICAiTWFyY2giLAogICAgIkFwcmlsIiwKICAgICJNYXkiLAogICAgIkp1bmUiLAogICAgIkp1bHkiLAogICAgIkF1Z3VzdCIsCiAgICAiU2VwdGVtYmVyIiwKICAgICJPY3RvYmVyIiwKICAgICJOb3ZlbWJlciIsCiAgICAiRGVjZW1iZXIiLAogIF07CiAgLy8gdHJhaWxpbmcgYmFja3NsYXNoIC0+IChkcm9wcGVkKQogIC8vIGEgYmFja3NsYXNoIGZvbGxvd2VkIGJ5IGFueSBjaGFyYWN0ZXIgKGluY2x1ZGluZyBiYWNrc2xhc2gpIC0+IHRoZSBjaGFyYWN0ZXIKICAvLyBlbXB0eSBzdHJpbmcgLT4gZW1wdHkgc3RyaW5nCiAgY29uc3QgZm9ybWF0Q2hyID0gL1xcPyguPykvZ2k7CiAgY29uc3QgZm9ybWF0Q2hyQ2IgPSBmdW5jdGlvbiAodCwgcykgewogICAgcmV0dXJuIGZbdF0gPyBmW3RdKCkgOiBzOwogIH07CiAgY29uc3QgX3BhZCA9IGZ1bmN0aW9uIChuLCBjKSB7CiAgICBuID0gU3RyaW5nKG4pOwogICAgd2hpbGUgKG4ubGVuZ3RoIDwgYykgewogICAgICBuID0gIjAiICsgbjsKICAgIH0KICAgIHJldHVybiBuOwogIH07CiAgZiA9IHsKICAgIC8vIERheQogICAgZCgpIHsKICAgICAgLy8gRGF5IG9mIG1vbnRoIHcvbGVhZGluZyAwOyAwMS4uMzEKICAgICAgcmV0dXJuIF9wYWQoZi5qKCksIDIpOwogICAgfSwKICAgIEQoKSB7CiAgICAgIC8vIFNob3J0aGFuZCBkYXkgbmFtZTsgTW9uLi4uU3VuCiAgICAgIHJldHVybiBmLmwoKS5zbGljZSgwLCAzKTsKICAgIH0sCiAgICBqKCkgewogICAgICAvLyBEYXkgb2YgbW9udGg7IDEuLjMxCiAgICAgIHJldHVybiBqc2RhdGUuZ2V0RGF0ZSgpOwogICAgfSwKICAgIGwoKSB7CiAgICAgIC8vIEZ1bGwgZGF5IG5hbWU7IE1vbmRheS4uLlN1bmRheQogICAgICByZXR1cm4gdHh0X3dvcmRzW2YudygpXSArICJkYXkiOwogICAgfSwKICAgIE4oKSB7CiAgICAgIC8vIElTTy04NjAxIGRheSBvZiB3ZWVrOyAxW01vbl0uLjdbU3VuXQogICAgICByZXR1cm4gZi53KCkgfHwgNzsKICAgIH0sCiAgICBTKCkgewogICAgICAvLyBPcmRpbmFsIHN1ZmZpeCBmb3IgZGF5IG9mIG1vbnRoOyBzdCwgbmQsIHJkLCB0aAogICAgICBjb25zdCBqID0gZi5qKCk7CiAgICAgIGxldCBpID0gaiAlIDEwOwogICAgICBpZiAoaSA8PSAzICYmIHBhcnNlSW50KChqICUgMTAwKSAvIDEwLCAxMCkgPT09IDEpIHsKICAgICAgICBpID0gMDsKICAgICAgfQogICAgICByZXR1cm4gWyJzdCIsICJuZCIsICJyZCJdW2kgLSAxXSB8fCAidGgiOwogICAgfSwKICAgIHcoKSB7CiAgICAgIC8vIERheSBvZiB3ZWVrOyAwW1N1bl0uLjZbU2F0XQogICAgICByZXR1cm4ganNkYXRlLmdldERheSgpOwogICAgfSwKICAgIHooKSB7CiAgICAgIC8vIERheSBvZiB5ZWFyOyAwLi4zNjUKICAgICAgY29uc3QgYSA9IG5ldyBEYXRlKGYuWSgpLCBmLm4oKSAtIDEsIGYuaigpKTsKICAgICAgY29uc3QgYiA9IG5ldyBEYXRlKGYuWSgpLCAwLCAxKTsKICAgICAgcmV0dXJuIE1hdGgucm91bmQoKGEgLSBiKSAvIDg2NGU1KTsKICAgIH0sCgogICAgLy8gV2VlawogICAgVygpIHsKICAgICAgLy8gSVNPLTg2MDEgd2VlayBudW1iZXIKICAgICAgY29uc3QgYSA9IG5ldyBEYXRlKGYuWSgpLCBmLm4oKSAtIDEsIGYuaigpIC0gZi5OKCkgKyAzKTsKICAgICAgY29uc3QgYiA9IG5ldyBEYXRlKGEuZ2V0RnVsbFllYXIoKSwgMCwgNCk7CiAgICAgIHJldHVybiBfcGFkKDEgKyBNYXRoLnJvdW5kKChhIC0gYikgLyA4NjRlNSAvIDcpLCAyKTsKICAgIH0sCgogICAgLy8gTW9udGgKICAgIEYoKSB7CiAgICAgIC8vIEZ1bGwgbW9udGggbmFtZTsgSmFudWFyeS4uLkRlY2VtYmVyCiAgICAgIHJldHVybiB0eHRfd29yZHNbNiArIGYubigpXTsKICAgIH0sCiAgICBtKCkgewogICAgICAvLyBNb250aCB3L2xlYWRpbmcgMDsgMDEuLi4xMgogICAgICByZXR1cm4gX3BhZChmLm4oKSwgMik7CiAgICB9LAogICAgTSgpIHsKICAgICAgLy8gU2hvcnRoYW5kIG1vbnRoIG5hbWU7IEphbi4uLkRlYwogICAgICByZXR1cm4gZi5GKCkuc2xpY2UoMCwgMyk7CiAgICB9LAogICAgbigpIHsKICAgICAgLy8gTW9udGg7IDEuLi4xMgogICAgICByZXR1cm4ganNkYXRlLmdldE1vbnRoKCkgKyAxOwogICAgfSwKICAgIHQoKSB7CiAgICAgIC8vIERheXMgaW4gbW9udGg7IDI4Li4uMzEKICAgICAgcmV0dXJuIG5ldyBEYXRlKGYuWSgpLCBmLm4oKSwgMCkuZ2V0RGF0ZSgpOwogICAgfSwKCiAgICAvLyBZZWFyCiAgICBMKCkgewogICAgICAvLyBJcyBsZWFwIHllYXI/OyAwIG9yIDEKICAgICAgY29uc3QgaiA9IGYuWSgpOwogICAgICByZXR1cm4gKChqICUgNCA9PT0gMCkgJiAoaiAlIDEwMCAhPT0gMCkpIHwgKGogJSA0MDAgPT09IDApOwogICAgfSwKICAgIG8oKSB7CiAgICAgIC8vIElTTy04NjAxIHllYXIKICAgICAgY29uc3QgbiA9IGYubigpOwogICAgICBjb25zdCBXID0gZi5XKCk7CiAgICAgIGNvbnN0IFkgPSBmLlkoKTsKICAgICAgLy8gZXNsaW50LWRpc2FibGUtbmV4dC1saW5lCiAgICAgIHJldHVybiBZICsgKG4gPT09IDEyICYmIFcgPCA5ID8gMSA6IG4gPT09IDEgJiYgVyA+IDkgPyAtMSA6IDApOwogICAgfSwKICAgIFkoKSB7CiAgICAgIC8vIEZ1bGwgeWVhcjsgZS5nLiAxOTgwLi4uMjAxMAogICAgICByZXR1cm4ganNkYXRlLmdldEZ1bGxZZWFyKCk7CiAgICB9LAogICAgeSgpIHsKICAgICAgLy8gTGFzdCB0d28gZGlnaXRzIG9mIHllYXI7IDAwLi4uOTkKICAgICAgcmV0dXJuIGYuWSgpLnRvU3RyaW5nKCkuc2xpY2UoLTIpOwogICAgfSwKCiAgICAvLyBUaW1lCiAgICBhKCkgewogICAgICAvLyBhbSBvciBwbQogICAgICByZXR1cm4ganNkYXRlLmdldEhvdXJzKCkgPiAxMSA/ICJwbSIgOiAiYW0iOwogICAgfSwKICAgIEEoKSB7CiAgICAgIC8vIEFNIG9yIFBNCiAgICAgIHJldHVybiBmLmEoKS50b1VwcGVyQ2FzZSgpOwogICAgfSwKICAgIEIoKSB7CiAgICAgIC8vIFN3YXRjaCBJbnRlcm5ldCB0aW1lOyAwMDAuLjk5OQogICAgICBjb25zdCBIID0ganNkYXRlLmdldFVUQ0hvdXJzKCkgKiAzNmUyOwogICAgICAvLyBIb3VycwogICAgICBjb25zdCBpID0ganNkYXRlLmdldFVUQ01pbnV0ZXMoKSAqIDYwOwogICAgICAvLyBNaW51dGVzCiAgICAgIGNvbnN0IHMgPSBqc2RhdGUuZ2V0VVRDU2Vjb25kcygpOyAvLyBTZWNvbmRzCiAgICAgIHJldHVybiBfcGFkKE1hdGguZmxvb3IoKEggKyBpICsgcyArIDM2ZTIpIC8gODYuNCkgJSAxZTMsIDMpOwogICAgfSwKICAgIGcoKSB7CiAgICAgIC8vIDEyLUhvdXJzOyAxLi4xMgogICAgICByZXR1cm4gZi5HKCkgJSAxMiB8fCAxMjsKICAgIH0sCiAgICBHKCkgewogICAgICAvLyAyNC1Ib3VyczsgMC4uMjMKICAgICAgcmV0dXJuIGpzZGF0ZS5nZXRIb3VycygpOwogICAgfSwKICAgIGgoKSB7CiAgICAgIC8vIDEyLUhvdXJzIHcvbGVhZGluZyAwOyAwMS4uMTIKICAgICAgcmV0dXJuIF9wYWQoZi5nKCksIDIpOwogICAgfSwKICAgIEgoKSB7CiAgICAgIC8vIDI0LUhvdXJzIHcvbGVhZGluZyAwOyAwMC4uMjMKICAgICAgcmV0dXJuIF9wYWQoZi5HKCksIDIpOwogICAgfSwKICAgIGkoKSB7CiAgICAgIC8vIE1pbnV0ZXMgdy9sZWFkaW5nIDA7IDAwLi41OQogICAgICByZXR1cm4gX3BhZChqc2RhdGUuZ2V0TWludXRlcygpLCAyKTsKICAgIH0sCiAgICBzKCkgewogICAgICAvLyBTZWNvbmRzIHcvbGVhZGluZyAwOyAwMC4uNTkKICAgICAgcmV0dXJuIF9wYWQoanNkYXRlLmdldFNlY29uZHMoKSwgMik7CiAgICB9LAogICAgdSgpIHsKICAgICAgLy8gTWljcm9zZWNvbmRzOyAwMDAwMDAtOTk5MDAwCiAgICAgIHJldHVybiBfcGFkKGpzZGF0ZS5nZXRNaWxsaXNlY29uZHMoKSAqIDEwMDAsIDYpOwogICAgfSwKCiAgICAvLyBUaW1lem9uZQogICAgZSgpIHsKICAgICAgLy8gVGltZXpvbmUgaWRlbnRpZmllcjsgZS5nLiBBdGxhbnRpYy9Bem9yZXMsIC4uLgogICAgICAvLyBUaGUgZm9sbG93aW5nIHdvcmtzLCBidXQgcmVxdWlyZXMgaW5jbHVzaW9uIG9mIHRoZSB2ZXJ5IGxhcmdlCiAgICAgIC8vIHRpbWV6b25lX2FiYnJldmlhdGlvbnNfbGlzdCgpIGZ1bmN0aW9uLgogICAgICAvKiAgICAgICAgICAgICAgcmV0dXJuIHRoYXQuZGF0ZV9kZWZhdWx0X3RpbWV6b25lX2dldCgpOwogICAgICAgKi8KICAgICAgdGhyb3cgbmV3IEVycm9yKAogICAgICAgICJOb3Qgc3VwcG9ydGVkIChzZWUgc291cmNlIGNvZGUgb2YgZGF0ZSgpIGZvciB0aW1lem9uZSBvbiBob3cgdG8gYWRkIHN1cHBvcnQpIgogICAgICApOwogICAgfSwKICAgIEkoKSB7CiAgICAgIC8vIERTVCBvYnNlcnZlZD87IDAgb3IgMQogICAgICAvLyBDb21wYXJlcyBKYW4gMSBtaW51cyBKYW4gMSBVVEMgdG8gSnVsIDEgbWludXMgSnVsIDEgVVRDLgogICAgICAvLyBJZiB0aGV5IGFyZSBub3QgZXF1YWwsIHRoZW4gRFNUIGlzIG9ic2VydmVkLgogICAgICBjb25zdCBhID0gbmV3IERhdGUoZi5ZKCksIDApOwogICAgICAvLyBKYW4gMQogICAgICBjb25zdCBjID0gRGF0ZS5VVEMoZi5ZKCksIDApOwogICAgICAvLyBKYW4gMSBVVEMKICAgICAgY29uc3QgYiA9IG5ldyBEYXRlKGYuWSgpLCA2KTsKICAgICAgLy8gSnVsIDEKICAgICAgY29uc3QgZCA9IERhdGUuVVRDKGYuWSgpLCA2KTsgLy8gSnVsIDEgVVRDCiAgICAgIHJldHVybiBhIC0gYyAhPT0gYiAtIGQgPyAxIDogMDsKICAgIH0sCiAgICBPKCkgewogICAgICAvLyBEaWZmZXJlbmNlIHRvIEdNVCBpbiBob3VyIGZvcm1hdDsgZS5nLiArMDIwMAogICAgICBjb25zdCB0em8gPSBqc2RhdGUuZ2V0VGltZXpvbmVPZmZzZXQoKTsKICAgICAgY29uc3QgYSA9IE1hdGguYWJzKHR6byk7CiAgICAgIHJldHVybiAoCiAgICAgICAgKHR6byA+IDAgPyAiLSIgOiAiKyIpICsgX3BhZChNYXRoLmZsb29yKGEgLyA2MCkgKiAxMDAgKyAoYSAlIDYwKSwgNCkKICAgICAgKTsKICAgIH0sCiAgICBQKCkgewogICAgICAvLyBEaWZmZXJlbmNlIHRvIEdNVCB3L2NvbG9uOyBlLmcuICswMjowMAogICAgICBjb25zdCBPID0gZi5PKCk7CiAgICAgIHJldHVybiBPLnN1YnN0cigwLCAzKSArICI6IiArIE8uc3Vic3RyKDMsIDIpOwogICAgfSwKICAgIFQoKSB7CiAgICAgIC8vIFRpbWV6b25lIGFiYnJldmlhdGlvbjsgZS5nLiBFU1QsIE1EVCwgLi4uCiAgICAgIC8vIFRoZSBmb2xsb3dpbmcgd29ya3MsIGJ1dCByZXF1aXJlcyBpbmNsdXNpb24gb2YgdGhlIHZlcnkKICAgICAgLy8gbGFyZ2UgdGltZXpvbmVfYWJicmV2aWF0aW9uc19saXN0KCkgZnVuY3Rpb24uCiAgICAgIC8qICAgICAgICAgICAgICB2YXIgYWJiciwgaSwgb3MsIF9kZWZhdWx0OwogICAgICBpZiAoIXRhbC5sZW5ndGgpIHsKICAgICAgICB0YWwgPSB0aGF0LnRpbWV6b25lX2FiYnJldmlhdGlvbnNfbGlzdCgpOwogICAgICB9CiAgICAgIGlmICh0aGF0LnBocF9qcyAmJiB0aGF0LnBocF9qcy5kZWZhdWx0X3RpbWV6b25lKSB7CiAgICAgICAgX2RlZmF1bHQgPSB0aGF0LnBocF9qcy5kZWZhdWx0X3RpbWV6b25lOwogICAgICAgIGZvciAoYWJiciBpbiB0YWwpIHsKICAgICAgICAgIGZvciAoaSA9IDA7IGkgPCB0YWxbYWJicl0ubGVuZ3RoOyBpKyspIHsKICAgICAgICAgICAgaWYgKHRhbFthYmJyXVtpXS50aW1lem9uZV9pZCA9PT0gX2RlZmF1bHQpIHsKICAgICAgICAgICAgICByZXR1cm4gYWJici50b1VwcGVyQ2FzZSgpOwogICAgICAgICAgICB9CiAgICAgICAgICB9CiAgICAgICAgfQogICAgICB9CiAgICAgIGZvciAoYWJiciBpbiB0YWwpIHsKICAgICAgICBmb3IgKGkgPSAwOyBpIDwgdGFsW2FiYnJdLmxlbmd0aDsgaSsrKSB7CiAgICAgICAgICBvcyA9IC1qc2RhdGUuZ2V0VGltZXpvbmVPZmZzZXQoKSAqIDYwOwogICAgICAgICAgaWYgKHRhbFthYmJyXVtpXS5vZmZzZXQgPT09IG9zKSB7CiAgICAgICAgICAgIHJldHVybiBhYmJyLnRvVXBwZXJDYXNlKCk7CiAgICAgICAgICB9CiAgICAgICAgfQogICAgICB9CiAgICAgICovCiAgICAgIHJldHVybiAiVVRDIjsKICAgIH0sCiAgICBaKCkgewogICAgICAvLyBUaW1lem9uZSBvZmZzZXQgaW4gc2Vjb25kcyAoLTQzMjAwLi4uNTA0MDApCiAgICAgIHJldHVybiAtanNkYXRlLmdldFRpbWV6b25lT2Zmc2V0KCkgKiA2MDsKICAgIH0sCgogICAgLy8gRnVsbCBEYXRlL1RpbWUKICAgIGMoKSB7CiAgICAgIC8vIElTTy04NjAxIGRhdGUuCiAgICAgIHJldHVybiAiWS1tLWRcXFRIOmk6c1AiLnJlcGxhY2UoZm9ybWF0Q2hyLCBmb3JtYXRDaHJDYik7CiAgICB9LAogICAgcigpIHsKICAgICAgLy8gUkZDIDI4MjIKICAgICAgcmV0dXJuICJELCBkIE0gWSBIOmk6cyBPIi5yZXBsYWNlKGZvcm1hdENociwgZm9ybWF0Q2hyQ2IpOwogICAgfSwKICAgIFUoKSB7CiAgICAgIC8vIFNlY29uZHMgc2luY2UgVU5JWCBlcG9jaAogICAgICByZXR1cm4gKGpzZGF0ZSAvIDEwMDApIHwgMDsKICAgIH0sCiAgfTsKICB0aGlzLmRhdGUgPSBmdW5jdGlvbiAoZm9ybWF0LCB0aW1lc3RhbXApIHsKICAgIC8vIGVzbGludC1kaXNhYmxlLW5leHQtbGluZQogICAganNkYXRlID0KICAgICAgdGltZXN0YW1wID09PSB1bmRlZmluZWQKICAgICAgICA/IG5ldyBEYXRlKCkgLy8gTm90IHByb3ZpZGVkCiAgICAgICAgOiB0aW1lc3RhbXAgaW5zdGFuY2VvZiBEYXRlCiAgICAgICAgPyBuZXcgRGF0ZSh0aW1lc3RhbXApIC8vIEpTIERhdGUoKQogICAgICAgIDogbmV3IERhdGUodGltZXN0YW1wICogMTAwMCk7IC8vIFVOSVggdGltZXN0YW1wIChhdXRvLWNvbnZlcnQgdG8gaW50KQogICAgcmV0dXJuIGZvcm1hdC5yZXBsYWNlKGZvcm1hdENociwgZm9ybWF0Q2hyQ2IpOwogIH07CiAgcmV0dXJuIHRoaXMuZGF0ZShmb3JtYXQsIHRpbWVzdGFtcCk7Cn0KCi8qIGJyb3VnaHQgZnJvbSBodHRwczovL2dpdGh1Yi5jb20vS3lsZUFNYXRoZXdzL2RlZXBtZXJnZSBhbmQgdW53cmFwcGVkIGZyb20gVU1EICovCmZ1bmN0aW9uIGRlZXBtZXJnZSh0YXJnZXQsIHNyYykgewogIHZhciBhcnJheSA9IEFycmF5LmlzQXJyYXkoc3JjKTsKICB2YXIgZHN0ID0gKGFycmF5ICYmIFtdKSB8fCB7fTsKCiAgaWYgKGFycmF5KSB7CiAgICB0YXJnZXQgPSB0YXJnZXQgfHwgW107CiAgICBkc3QgPSBkc3QuY29uY2F0KHRhcmdldCk7CiAgICBzcmMuZm9yRWFjaChmdW5jdGlvbiAoZSwgaSkgewogICAgICBpZiAodHlwZW9mIGRzdFtpXSA9PT0gInVuZGVmaW5lZCIpIHsKICAgICAgICBkc3RbaV0gPSBlOwogICAgICB9IGVsc2UgaWYgKHR5cGVvZiBlID09PSAib2JqZWN0IikgewogICAgICAgIGRzdFtpXSA9IGRlZXBtZXJnZSh0YXJnZXRbaV0sIGUpOwogICAgICB9IGVsc2UgewogICAgICAgIGlmICh0YXJnZXQuaW5kZXhPZihlKSA9PT0gLTEpIHsKICAgICAgICAgIGRzdC5wdXNoKGUpOwogICAgICAgIH0KICAgICAgfQogICAgfSk7CiAgfSBlbHNlIHsKICAgIGlmICh0YXJnZXQgJiYgdHlwZW9mIHRhcmdldCA9PT0gIm9iamVjdCIpIHsKICAgICAgT2JqZWN0LmtleXModGFyZ2V0KS5mb3JFYWNoKGZ1bmN0aW9uIChrZXkpIHsKICAgICAgICBkc3Rba2V5XSA9IHRhcmdldFtrZXldOwogICAgICB9KTsKICAgIH0KICAgIE9iamVjdC5rZXlzKHNyYykuZm9yRWFjaChmdW5jdGlvbiAoa2V5KSB7CiAgICAgIGlmICh0eXBlb2Ygc3JjW2tleV0gIT09ICJvYmplY3QiIHx8ICFzcmNba2V5XSkgewogICAgICAgIGRzdFtrZXldID0gc3JjW2tleV07CiAgICAgIH0gZWxzZSB7CiAgICAgICAgaWYgKCF0YXJnZXRba2V5XSkgewogICAgICAgICAgZHN0W2tleV0gPSBzcmNba2V5XTsKICAgICAgICB9IGVsc2UgewogICAgICAgICAgZHN0W2tleV0gPSBkZWVwbWVyZ2UodGFyZ2V0W2tleV0sIHNyY1trZXldKTsKICAgICAgICB9CiAgICAgIH0KICAgIH0pOwogIH0KCiAgcmV0dXJuIGRzdDsKfQo="})

	// builtin_log.js
	modules = append(modules, JsModule{File: "builtin_log.js", Code: "ZnVuY3Rpb24gcHJpbnRsbiguLi5hcmdzKSB7CiAgcHJpbnQoLi4uYXJncyk7CiAgcHJpbnQoIlxuIik7Cn0KCmNvbnN0IGxvZyA9IHt9Owp7CiAgY29uc3QgbGV2ZWxzID0geyBFUlJPUjogMSwgSU5GTzogMiwgREVCVUc6IDMgfTsKICBjb25zdCBsb2dMZXZlbCA9IGxldmVsc1soX19lbnZbIkpTU0hfTE9HIl0gfHwgIklORk8iKS50b1VwcGVyQ2FzZSgpXTsKCiAgY29uc3QgcmVzZXQgPSBgXHUwMDFiWzBtYDsKCiAgZnVuY3Rpb24gcmVkKGxpbmUpIHsKICAgIHJldHVybiBgXHUwMDFiWzMxOzFtJHtsaW5lfSR7cmVzZXR9YDsKICB9CgogIGZ1bmN0aW9uIGdyZWVuKGxpbmUpIHsKICAgIHJldHVybiBgXHUwMDFiWzMyOzFtJHtsaW5lfSR7cmVzZXR9YDsKICB9CgogIGZ1bmN0aW9uIGdyYXkobGluZSkgewogICAgcmV0dXJuIGBcdTAwMWJbMjsxbSR7bGluZX0ke3Jlc2V0fWA7CiAgfQoKICBsb2cuZGVidWcgPSBmdW5jdGlvbiBkZWJ1ZyhtZXNzYWdlLCAuLi5hcmdzKSB7CiAgICBpZiAobG9nTGV2ZWwgPj0gbGV2ZWxzLkRFQlVHKSB7CiAgICAgIHN0ZG91dGxvZyhncmF5KGZvcm1hdChtZXNzYWdlLCAuLi5hcmdzKSkpOwogICAgfQogIH07CgogIGxvZy5pbmZvID0gZnVuY3Rpb24gaW5mbyhtZXNzYWdlLCAuLi5hcmdzKSB7CiAgICBpZiAobG9nTGV2ZWwgPj0gbGV2ZWxzLklORk8pIHsKICAgICAgc3Rkb3V0bG9nKGdyZWVuKGZvcm1hdChtZXNzYWdlLCAuLi5hcmdzKSkpOwogICAgfQogIH07CgogIGxvZy5lcnJvciA9IGZ1bmN0aW9uIGVycm9yKG1lc3NhZ2UsIC4uLmFyZ3MpIHsKICAgIGlmIChsb2dMZXZlbCA+PSBsZXZlbHMuRVJST1IpIHsKICAgICAgc3RkZXJybG9nKHJlZChmb3JtYXQobWVzc2FnZSwgLi4uYXJncykpKTsKICAgIH0KICB9OwoKICBsb2cuZmF0YWwgPSBmdW5jdGlvbiBmYXRhbChtZXNzYWdlLCAuLi5hcmdzKSB7CiAgICBsb2cuZXJyb3IobWVzc2FnZSwgLi4uYXJncyk7CiAgICBleGl0KDEpOwogIH07Cn0K"})
}
