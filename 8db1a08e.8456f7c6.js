(window.webpackJsonp=window.webpackJsonp||[]).push([[10],{139:function(e,t,n){"use strict";n.r(t),t.default=n.p+"assets/images/rpc-110606ff9ff87661d7006a8901018c91.png"},140:function(e,t,n){"use strict";n.r(t),t.default=n.p+"assets/images/rpcedit-5c5726607877597a8cee0b57ec482935.png"},80:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return o})),n.d(t,"metadata",(function(){return i})),n.d(t,"toc",(function(){return s})),n.d(t,"default",(function(){return l}));var r=n(3),a=n(7),c=(n(0),n(95)),o={id:"rpcs",title:"RPCs"},i={unversionedId:"overview/rpcs",id:"overview/rpcs",isDocsHomePage:!1,title:"RPCs",description:"Aft's API can cover most ordinary client needs. But sometimes you just need an escape hatch; for that, Aft includes a scriptable RPC system.",source:"@site/docs/overview/rpcs.md",slug:"/overview/rpcs",permalink:"/docs/overview/rpcs",version:"current",sidebar:"main",previous:{title:"API",permalink:"/docs/overview/api"},next:{title:"Access",permalink:"/docs/overview/access"}},s=[{value:"Starlark",id:"starlark",children:[]}],p={toc:s};function l(e){var t=e.components,o=Object(a.a)(e,["components"]);return Object(c.b)("wrapper",Object(r.a)({},p,o,{components:t,mdxType:"MDXLayout"}),Object(c.b)("p",null,"Aft's API can cover most ordinary client needs. But sometimes you just need an escape hatch; for that, Aft includes a scriptable RPC system."),Object(c.b)("p",null,Object(c.b)("img",{alt:"Screenshot of the rpc page",src:n(139).default})),Object(c.b)("p",null,"RPCs can be written in Starlark. They are passed two arguments; first, a handle to a Starlark version of the Aft API and second, a dictionary that is a decoded version of a JSON object passed from the client."),Object(c.b)("p",null,"The RPCs are exposed in the following URL format:"),Object(c.b)("pre",null,Object(c.b)("code",Object(r.a)({parentName:"pre"},{}),"https://$BASE_URL/api/rpc.$RPC_NAME\n")),Object(c.b)("p",null,'The RPC endpoint accepts a JSON object with a single key, "args":'),Object(c.b)("pre",null,Object(c.b)("code",Object(r.a)({parentName:"pre"},{}),'{\n    "args": {\n        "foo": "bar"\n    }\n}\n')),Object(c.b)("h2",{id:"starlark"},"Starlark"),Object(c.b)("p",null,Object(c.b)("img",{alt:"Screenshot of the rpc edit page",src:n(140).default})),Object(c.b)("p",null,'To write an RPC in Starlark, write a script with a function, "main," of two arguments: ',Object(c.b)("inlineCode",{parentName:"p"},"aft"),", a handle to the API and authentication methods, and ",Object(c.b)("inlineCode",{parentName:"p"},"data"),", a single json object sent by the client."),Object(c.b)("p",null,"Aft's API methods can be accessed by calling them on the ",Object(c.b)("inlineCode",{parentName:"p"},"api")," object like so:"),Object(c.b)("pre",null,Object(c.b)("code",Object(r.a)({parentName:"pre"},{className:"language-python"}),'def main(aft, data):\n    user = aft.api.findOne("users", {"where": {"name": "Andrew"}})\n    return user.name  # returns "Andrew"\n')),Object(c.b)("p",null,"Modifying records returned by the ",Object(c.b)("inlineCode",{parentName:"p"},"aft.api")," methods will not have an effect on the datastore. To mutate the datastore, use the mutation api calls. "),Object(c.b)("pre",null,Object(c.b)("code",Object(r.a)({parentName:"pre"},{className:"language-python"}),'def main(aft, data):\n    user = aft.api.findOne("users", {"where": {"name": "Andrew"}})\n    aft.api.update("users", {"where": {"id": user.id}, "data": {"name": "Werdna"}})\n')))}l.isMDXComponent=!0},95:function(e,t,n){"use strict";n.d(t,"a",(function(){return u})),n.d(t,"b",(function(){return f}));var r=n(0),a=n.n(r);function c(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function o(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?o(Object(n),!0).forEach((function(t){c(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):o(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},c=Object.keys(e);for(r=0;r<c.length;r++)n=c[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var c=Object.getOwnPropertySymbols(e);for(r=0;r<c.length;r++)n=c[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var p=a.a.createContext({}),l=function(e){var t=a.a.useContext(p),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},u=function(e){var t=l(e.components);return a.a.createElement(p.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},b=a.a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,c=e.originalType,o=e.parentName,p=s(e,["components","mdxType","originalType","parentName"]),u=l(n),b=r,f=u["".concat(o,".").concat(b)]||u[b]||d[b]||c;return n?a.a.createElement(f,i(i({ref:t},p),{},{components:n})):a.a.createElement(f,i({ref:t},p))}));function f(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var c=n.length,o=new Array(c);o[0]=b;var i={};for(var s in t)hasOwnProperty.call(t,s)&&(i[s]=t[s]);i.originalType=e,i.mdxType="string"==typeof e?e:r,o[1]=i;for(var p=2;p<c;p++)o[p]=n[p];return a.a.createElement.apply(null,o)}return a.a.createElement.apply(null,n)}b.displayName="MDXCreateElement"}}]);