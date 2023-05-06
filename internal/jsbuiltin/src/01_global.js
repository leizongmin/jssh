function randomstring(size, chars) {
  size = size || 6;
  chars =
    chars || "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const max = chars.length;
  let ret = "";
  for (let i = 0; i < size; i++) {
    ret += chars.charAt(Math.floor(Math.random() * max));
  }
  return ret;
}

function formatdate(format, timestamp) {
  //  discuss at: http://phpjs.org/functions/date/
  //   example 1: date('H:m:s \\m \\i\\s \\m\\o\\n\\t\\h', 1062402400);
  //   returns 1: '09:09:40 m is month'
  //   example 2: date('F j, Y, g:i a', 1062462400);
  //   returns 2: 'September 2, 2003, 2:26 am'
  //   example 3: date('Y W o', 1062462400);
  //   returns 3: '2003 36 2003'
  //   example 4: x = date('Y m d', (new Date()).getTime()/1000);
  //   example 4: (x+'').length == 10 // 2009 01 09
  //   returns 4: true
  //   example 5: date('W', 1104534000);
  //   returns 5: '53'
  //   example 6: date('B t', 1104534000);
  //   returns 6: '999 31'
  //   example 7: date('W U', 1293750000.82); // 2010-12-31
  //   returns 7: '52 1293750000'
  //   example 8: date('W', 1293836400); // 2011-01-01
  //   returns 8: '52'
  //   example 9: date('W Y-m-d', 1293974054); // 2011-01-02
  //   returns 9: '52 2011-01-02'

  let jsdate, f;
  // Keep this here (works, but for code commented-out below for file size reasons)
  // var tal= [];
  const txt_words = [
    "Sun",
    "Mon",
    "Tues",
    "Wednes",
    "Thurs",
    "Fri",
    "Satur",
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December",
  ];
  // trailing backslash -> (dropped)
  // a backslash followed by any character (including backslash) -> the character
  // empty string -> empty string
  const formatChr = /\\?(.?)/gi;
  const formatChrCb = function (t, s) {
    return f[t] ? f[t]() : s;
  };
  const _pad = function (n, c) {
    n = String(n);
    while (n.length < c) {
      n = "0" + n;
    }
    return n;
  };
  f = {
    // Day
    d() {
      // Day of month w/leading 0; 01..31
      return _pad(f.j(), 2);
    },
    D() {
      // Shorthand day name; Mon...Sun
      return f.l().slice(0, 3);
    },
    j() {
      // Day of month; 1..31
      return jsdate.getDate();
    },
    l() {
      // Full day name; Monday...Sunday
      return txt_words[f.w()] + "day";
    },
    N() {
      // ISO-8601 day of week; 1[Mon]..7[Sun]
      return f.w() || 7;
    },
    S() {
      // Ordinal suffix for day of month; st, nd, rd, th
      const j = f.j();
      let i = j % 10;
      if (i <= 3 && parseInt((j % 100) / 10, 10) === 1) {
        i = 0;
      }
      return ["st", "nd", "rd"][i - 1] || "th";
    },
    w() {
      // Day of week; 0[Sun]..6[Sat]
      return jsdate.getDay();
    },
    z() {
      // Day of year; 0..365
      const a = new Date(f.Y(), f.n() - 1, f.j());
      const b = new Date(f.Y(), 0, 1);
      return Math.round((a - b) / 864e5);
    },

    // Week
    W() {
      // ISO-8601 week number
      const a = new Date(f.Y(), f.n() - 1, f.j() - f.N() + 3);
      const b = new Date(a.getFullYear(), 0, 4);
      return _pad(1 + Math.round((a - b) / 864e5 / 7), 2);
    },

    // Month
    F() {
      // Full month name; January...December
      return txt_words[6 + f.n()];
    },
    m() {
      // Month w/leading 0; 01...12
      return _pad(f.n(), 2);
    },
    M() {
      // Shorthand month name; Jan...Dec
      return f.F().slice(0, 3);
    },
    n() {
      // Month; 1...12
      return jsdate.getMonth() + 1;
    },
    t() {
      // Days in month; 28...31
      return new Date(f.Y(), f.n(), 0).getDate();
    },

    // Year
    L() {
      // Is leap year?; 0 or 1
      const j = f.Y();
      return ((j % 4 === 0) & (j % 100 !== 0)) | (j % 400 === 0);
    },
    o() {
      // ISO-8601 year
      const n = f.n();
      const W = f.W();
      const Y = f.Y();
      // eslint-disable-next-line
      return Y + (n === 12 && W < 9 ? 1 : n === 1 && W > 9 ? -1 : 0);
    },
    Y() {
      // Full year; e.g. 1980...2010
      return jsdate.getFullYear();
    },
    y() {
      // Last two digits of year; 00...99
      return f.Y().toString().slice(-2);
    },

    // Time
    a() {
      // am or pm
      return jsdate.getHours() > 11 ? "pm" : "am";
    },
    A() {
      // AM or PM
      return f.a().toUpperCase();
    },
    B() {
      // Swatch Internet time; 000..999
      const H = jsdate.getUTCHours() * 36e2;
      // Hours
      const i = jsdate.getUTCMinutes() * 60;
      // Minutes
      const s = jsdate.getUTCSeconds(); // Seconds
      return _pad(Math.floor((H + i + s + 36e2) / 86.4) % 1e3, 3);
    },
    g() {
      // 12-Hours; 1..12
      return f.G() % 12 || 12;
    },
    G() {
      // 24-Hours; 0..23
      return jsdate.getHours();
    },
    h() {
      // 12-Hours w/leading 0; 01..12
      return _pad(f.g(), 2);
    },
    H() {
      // 24-Hours w/leading 0; 00..23
      return _pad(f.G(), 2);
    },
    i() {
      // Minutes w/leading 0; 00..59
      return _pad(jsdate.getMinutes(), 2);
    },
    s() {
      // Seconds w/leading 0; 00..59
      return _pad(jsdate.getSeconds(), 2);
    },
    u() {
      // Microseconds; 000000-999000
      return _pad(jsdate.getMilliseconds() * 1000, 6);
    },

    // Timezone
    e() {
      // Timezone identifier; e.g. Atlantic/Azores, ...
      // The following works, but requires inclusion of the very large
      // timezone_abbreviations_list() function.
      /*              return that.date_default_timezone_get();
       */
      throw new Error(
        "Not supported (see source code of date() for timezone on how to add support)"
      );
    },
    I() {
      // DST observed?; 0 or 1
      // Compares Jan 1 minus Jan 1 UTC to Jul 1 minus Jul 1 UTC.
      // If they are not equal, then DST is observed.
      const a = new Date(f.Y(), 0);
      // Jan 1
      const c = Date.UTC(f.Y(), 0);
      // Jan 1 UTC
      const b = new Date(f.Y(), 6);
      // Jul 1
      const d = Date.UTC(f.Y(), 6); // Jul 1 UTC
      return a - c !== b - d ? 1 : 0;
    },
    O() {
      // Difference to GMT in hour format; e.g. +0200
      const tzo = jsdate.getTimezoneOffset();
      const a = Math.abs(tzo);
      return (
        (tzo > 0 ? "-" : "+") + _pad(Math.floor(a / 60) * 100 + (a % 60), 4)
      );
    },
    P() {
      // Difference to GMT w/colon; e.g. +02:00
      const O = f.O();
      return O.substr(0, 3) + ":" + O.substr(3, 2);
    },
    T() {
      // Timezone abbreviation; e.g. EST, MDT, ...
      // The following works, but requires inclusion of the very
      // large timezone_abbreviations_list() function.
      /*              var abbr, i, os, _default;
      if (!tal.length) {
        tal = that.timezone_abbreviations_list();
      }
      if (that.php_js && that.php_js.default_timezone) {
        _default = that.php_js.default_timezone;
        for (abbr in tal) {
          for (i = 0; i < tal[abbr].length; i++) {
            if (tal[abbr][i].timezone_id === _default) {
              return abbr.toUpperCase();
            }
          }
        }
      }
      for (abbr in tal) {
        for (i = 0; i < tal[abbr].length; i++) {
          os = -jsdate.getTimezoneOffset() * 60;
          if (tal[abbr][i].offset === os) {
            return abbr.toUpperCase();
          }
        }
      }
      */
      return "UTC";
    },
    Z() {
      // Timezone offset in seconds (-43200...50400)
      return -jsdate.getTimezoneOffset() * 60;
    },

    // Full Date/Time
    c() {
      // ISO-8601 date.
      return "Y-m-d\\TH:i:sP".replace(formatChr, formatChrCb);
    },
    r() {
      // RFC 2822
      return "D, d M Y H:i:s O".replace(formatChr, formatChrCb);
    },
    U() {
      // Seconds since UNIX epoch
      return (jsdate / 1000) | 0;
    },
  };
  this.date = function (format, timestamp) {
    // eslint-disable-next-line
    jsdate =
      timestamp === undefined
        ? new Date() // Not provided
        : timestamp instanceof Date
        ? new Date(timestamp) // JS Date()
        : new Date(timestamp * 1000); // UNIX timestamp (auto-convert to int)
    return format.replace(formatChr, formatChrCb);
  };
  return this.date(format, timestamp);
}

/* brought from https://github.com/KyleAMathews/deepmerge and unwrapped from UMD */
function deepmerge(target, src) {
  var array = Array.isArray(src);
  var dst = (array && []) || {};

  if (array) {
    target = target || [];
    dst = dst.concat(target);
    src.forEach(function (e, i) {
      if (typeof dst[i] === "undefined") {
        dst[i] = e;
      } else if (typeof e === "object") {
        dst[i] = deepmerge(target[i], e);
      } else {
        if (target.indexOf(e) === -1) {
          dst.push(e);
        }
      }
    });
  } else {
    if (target && typeof target === "object") {
      Object.keys(target).forEach(function (key) {
        dst[key] = target[key];
      });
    }
    Object.keys(src).forEach(function (key) {
      if (typeof src[key] !== "object" || !src[key]) {
        dst[key] = src[key];
      } else {
        if (!target[key]) {
          dst[key] = src[key];
        } else {
          dst[key] = deepmerge(target[key], src[key]);
        }
      }
    });
  }

  return dst;
}
