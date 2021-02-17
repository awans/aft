(window.webpackJsonp=window.webpackJsonp||[]).push([[20],{142:function(e,t,r){"use strict";r.r(t),t.default=r.p+"assets/images/roles-225cd4309572409938e904b5174a7bef.png"},90:function(e,t,r){"use strict";r.r(t),r.d(t,"frontMatter",(function(){return o})),r.d(t,"metadata",(function(){return c})),r.d(t,"toc",(function(){return s})),r.d(t,"default",(function(){return l}));var n=r(3),a=(r(0),r(95));const o={id:"access",title:"Access"},c={unversionedId:"overview/access",id:"overview/access",isDocsHomePage:!1,title:"Access",description:"Screenshot of the roles page",source:"@site/docs/overview/access.md",slug:"/overview/access",permalink:"/docs/overview/access",version:"current",sidebar:"main",previous:{title:"RPCs",permalink:"/docs/overview/rpcs"},next:{title:"Identity",permalink:"/docs/overview/identity"}},s=[],i={toc:s};function l({components:e,...t}){return Object(a.b)("wrapper",Object(n.a)({},i,t,{components:e,mdxType:"MDXLayout"}),Object(a.b)("p",null,Object(a.b)("img",{alt:"Screenshot of the roles page",src:r(142).default})),Object(a.b)("p",null,"Aft's powerful API is paired with an equally powerful set of access controls. Aft is closed by default\u2014users with a given role can only access data explicitly granted by a policy."),Object(a.b)("p",null,'Policies are expressed on a per-interface and per-operation basis as "where" clauses as used in findMany queries. For example, to restrict access to just users named Andrew, one might have the following read-policy.'),Object(a.b)("pre",null,Object(a.b)("code",Object(n.a)({parentName:"pre"},{className:"language-json"}),'{\n    "name":"Andrew"\n}\n')),Object(a.b)("p",null,"They're also able to perform template string substitution to restrict on the basis of the current user ID. So to restrict a user to access only their own user data:"),Object(a.b)("pre",null,Object(a.b)("code",Object(n.a)({parentName:"pre"},{className:"language-json"}),'{\n    "id":"$USER_ID"\n}\n')),Object(a.b)("p",null,"Or for a model that had a relationship to user:"),Object(a.b)("pre",null,Object(a.b)("code",Object(n.a)({parentName:"pre"},{className:"language-json"}),'{\n    "user": {"id":"$USER_ID"}\n}\n')),Object(a.b)("p",null,"Connections and disconnections are allowed if and only if update is allowed on both records. Similarly, a user must retain update rights to any record that they update after the update is applied."))}l.isMDXComponent=!0},95:function(e,t,r){"use strict";r.d(t,"a",(function(){return p})),r.d(t,"b",(function(){return b}));var n=r(0),a=r.n(n);function o(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function c(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function s(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?c(Object(r),!0).forEach((function(t){o(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):c(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function i(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},o=Object.keys(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(n=0;n<o.length;n++)r=o[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var l=a.a.createContext({}),u=function(e){var t=a.a.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):s(s({},t),e)),r},p=function(e){var t=u(e.components);return a.a.createElement(l.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},f=a.a.forwardRef((function(e,t){var r=e.components,n=e.mdxType,o=e.originalType,c=e.parentName,l=i(e,["components","mdxType","originalType","parentName"]),p=u(r),f=n,b=p["".concat(c,".").concat(f)]||p[f]||d[f]||o;return r?a.a.createElement(b,s(s({ref:t},l),{},{components:r})):a.a.createElement(b,s({ref:t},l))}));function b(e,t){var r=arguments,n=t&&t.mdxType;if("string"==typeof e||n){var o=r.length,c=new Array(o);c[0]=f;var s={};for(var i in t)hasOwnProperty.call(t,i)&&(s[i]=t[i]);s.originalType=e,s.mdxType="string"==typeof e?e:n,c[1]=s;for(var l=2;l<o;l++)c[l]=r[l];return a.a.createElement.apply(null,c)}return a.a.createElement.apply(null,r)}f.displayName="MDXCreateElement"}}]);