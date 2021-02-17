(window.webpackJsonp=window.webpackJsonp||[]).push([[22],{92:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return a})),n.d(t,"metadata",(function(){return p})),n.d(t,"toc",(function(){return i})),n.d(t,"default",(function(){return c}));var r=n(3),o=(n(0),n(95));const a={id:"app-setup",title:"App Setup"},p={unversionedId:"tutorial/app-setup",id:"tutorial/app-setup",isDocsHomePage:!1,title:"App Setup",description:"First, let's make a new file, app.js, and move our App component over, making sure to export it.",source:"@site/docs/tutorial/03-app-setup.md",slug:"/tutorial/app-setup",permalink:"/docs/tutorial/app-setup",version:"current",sidebar:"main",previous:{title:"Frontend Setup",permalink:"/docs/tutorial/frontend-setup"},next:{title:"Login",permalink:"/docs/tutorial/login"}},i=[{value:"API client",id:"api-client",children:[]}],l={toc:i};function c({components:e,...t}){return Object(o.b)("wrapper",Object(r.a)({},l,t,{components:e,mdxType:"MDXLayout"}),Object(o.b)("p",null,"First, let's make a new file, ",Object(o.b)("inlineCode",{parentName:"p"},"app.js"),", and move our App component over, making sure to export it."),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-js",metastring:'title="app.js"',title:'"app.js"'}),"import {html} from 'https://unpkg.com/htm/preact/standalone.module.js'\n\nexport function App() {\n    return html`<h1>Hello Aft!</h1>`\n}\n")),Object(o.b)("p",null,"And then import it in our ",Object(o.b)("inlineCode",{parentName:"p"},"index.html")," file:"),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-html",metastring:'title="index.html"',title:'"index.html"'}),"<head>\n    <link rel=stylesheet href=\"./styles.css\" />\n    <script type=module>\n        import {html, render} from 'https://unpkg.com/htm/preact/standalone.module.js'\n        import {App} from './app.js'\n\n        render(html`<${App} />`, document.body);\n    <\/script>\n</head>\n")),Object(o.b)("p",null,"Hit refresh on the client\u2014you should still see your app rendering its greeting."),Object(o.b)("h2",{id:"api-client"},"API client"),Object(o.b)("p",null,"Now we'll add a small API client\u2014some objects that will make it easy for us to talk to Aft."),Object(o.b)("p",null,"Make a new file, ",Object(o.b)("inlineCode",{parentName:"p"},"aft.js"),", and add the following."),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-js",metastring:'title="aft.js"',title:'"aft.js"'}),"async function call(path, body) {\n    const result = await fetch(path, {\n        method: 'POST',\n        headers: {'Content-Type': 'application/json'},\n        body: JSON.stringify(body || {}),\n    })\n    const response = await result.json();\n    if (response.code) {\n        throw new Error(response.message);\n    }\n    return response.data;\n}\n\nconst curryProxy = (inner) => {\n    return new Proxy({}, {\n        get(_, prop)  { \n            return inner(prop) \n        }\n    })\n}\n\nexport default {\n    api: curryProxy((interfaceName) => curryProxy((method) => (params) => {\n        return call(\"api/\" + interfaceName + \".\" + method, params)\n    })),\n    rpc: curryProxy((rpcName) => (args) => {\n        return call(\"rpc/\" + rpcName, args)\n    }),\n}\n")),Object(o.b)("p",null,"The use of Proxy isn't really necessary, but it gives us a nice looking syntax for making API calls or RPCs to Aft. This short snippet is all you'll need in your app to use every bit of functionality Aft has to offer."),Object(o.b)("p",null,"Okay, nice work! Up next, we'll make our login UI and sign in to our app."))}c.isMDXComponent=!0},95:function(e,t,n){"use strict";n.d(t,"a",(function(){return u})),n.d(t,"b",(function(){return f}));var r=n(0),o=n.n(r);function a(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function p(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?p(Object(n),!0).forEach((function(t){a(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):p(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},a=Object.keys(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var c=o.a.createContext({}),s=function(e){var t=o.a.useContext(c),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},u=function(e){var t=s(e.components);return o.a.createElement(c.Provider,{value:t},e.children)},m={inlineCode:"code",wrapper:function(e){var t=e.children;return o.a.createElement(o.a.Fragment,{},t)}},d=o.a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,a=e.originalType,p=e.parentName,c=l(e,["components","mdxType","originalType","parentName"]),u=s(n),d=r,f=u["".concat(p,".").concat(d)]||u[d]||m[d]||a;return n?o.a.createElement(f,i(i({ref:t},c),{},{components:n})):o.a.createElement(f,i({ref:t},c))}));function f(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var a=n.length,p=new Array(a);p[0]=d;var i={};for(var l in t)hasOwnProperty.call(t,l)&&(i[l]=t[l]);i.originalType=e,i.mdxType="string"==typeof e?e:r,p[1]=i;for(var c=2;c<a;c++)p[c]=n[c];return o.a.createElement.apply(null,p)}return o.a.createElement.apply(null,n)}d.displayName="MDXCreateElement"}}]);