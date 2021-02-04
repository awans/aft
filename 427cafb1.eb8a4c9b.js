(window.webpackJsonp=window.webpackJsonp||[]).push([[4],{74:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return s})),n.d(t,"metadata",(function(){return i})),n.d(t,"toc",(function(){return l})),n.d(t,"default",(function(){return c}));var r=n(3),a=n(7),o=(n(0),n(95)),s={id:"login",title:"Login"},i={unversionedId:"tutorial/login",id:"tutorial/login",isDocsHomePage:!1,title:"Login",description:"Let's jump back and build our login widget back in app.js.",source:"@site/docs/tutorial/04-login.md",slug:"/tutorial/login",permalink:"/aft/docs/tutorial/login",version:"current",sidebar:"main",previous:{title:"App Setup",permalink:"/aft/docs/tutorial/app-setup"},next:{title:"User",permalink:"/aft/docs/tutorial/user"}},l=[{value:"Adding a user",id:"adding-a-user",children:[]}],u={toc:l};function c(e){var t=e.components,n=Object(a.a)(e,["components"]);return Object(o.b)("wrapper",Object(r.a)({},u,n,{components:t,mdxType:"MDXLayout"}),Object(o.b)("p",null,"Let's jump back and build our login widget back in ",Object(o.b)("inlineCode",{parentName:"p"},"app.js"),"."),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-js",metastring:'title="app.js"',title:'"app.js"'}),'import {html, useState} from \'https://unpkg.com/htm/preact/standalone.module.js\'\nimport aft from \'./aft.js\'\n\n\nexport function App (props) {\n    const [user, setUser] = useState(null);\n    return html`<${Login} setUser=${setUser} />`\n}\n\nfunction Login({setUser}) {\n    const [errorMessage, setErrorMessage] = useState(null);\n    const [email, setEmail] = useState("");\n    const [password, setPassword] = useState("");\n\n    return html`\n    <div class="box stack">\n        <input type=email placeholder="Email" \n            value=${email} \n            onInput=${(e) => setEmail(e.target.value)}/>\n        <input type=password placeholder="Password" \n            value=${password} \n            onInput=${(e) => setPassword(e.target.value)}/>\n        ${errorMessage && html`<div class=error>${errorMessage}</div>`}\n        <button>Sign in</button>\n    </div>`\n}\n')),Object(o.b)("p",null,"Refresh and take a look at our login box. Looking good!"),Object(o.b)("p",null,"Now we'll try and actually connect it to Aft. Add a ",Object(o.b)("inlineCode",{parentName:"p"},"submit")," callback and connect it to the button's ",Object(o.b)("inlineCode",{parentName:"p"},"onClick")," property."),Object(o.b)("p",null,"Notice how we're invoking the aft \"login\" RPC. Aft RPCs accept and return a single JSON object. We're able to just call the RPC by name like a native function thanks to the Proxy magic we did earlier in ",Object(o.b)("inlineCode",{parentName:"p"},"aft.js"),"."),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-js",metastring:'title="app.js"',title:'"app.js"'}),'function Login({setUser}) {\n    const [errorMessage, setErrorMessage] = useState(null);\n    const [email, setEmail] = useState("");\n    const [password, setPassword] = useState("");\n\n    const submit = useCallback(async () => {\n        setErrorMessage(null);\n        try {\n            const user = await aft.rpc.login({\n                email: email,\n                password: password,\n            });\n            setUser(user);\n        } catch (e) {\n            setErrorMessage(e.message);\n        }\n    }, [email, password])\n\n    return html`\n    <div class="box stack">\n        <input type=email placeholder="Email" \n            value=${email} \n            onInput=${(e) => setEmail(e.target.value)}/>\n        <input type=password placeholder="Password" \n            value=${password} \n            onInput=${(e) => setPassword(e.target.value)}/>\n        ${errorMessage && html`<div class=error>${errorMessage}</div>`}\n        <button onClick=${submit}>Sign in</button>\n    </div>`\n}\n')),Object(o.b)("h2",{id:"adding-a-user"},"Adding a user"),Object(o.b)("p",null,"If you go ahead and try to sign in to our tutorial app, you should of course get an error about the login not working\u2014we haven't added any users yet!"),Object(o.b)("p",null,"Open up Aft, and navigate to the ",Object(o.b)("strong",{parentName:"p"},"Terminal"),", and we'll create a user."),Object(o.b)("pre",null,Object(o.b)("code",Object(r.a)({parentName:"pre"},{className:"language-python"}),'def main(aft):\n    return aft.api.create("user", {"data": {\n                "email": "user@example.com", \n                "password": "coolpass",\n            }})\n')),Object(o.b)("p",null,"Press ",Object(o.b)("strong",{parentName:"p"},"Run"),", and you should see a JSON representation of the user just created, though the password is salted and hashed. "),Object(o.b)("p",null,"Go back to the tutorial app, and try signing in with your new credentials."),Object(o.b)("p",null,"In the next section, we'll finish up the login system, using the ",Object(o.b)("inlineCode",{parentName:"p"},"me")," RPC."))}c.isMDXComponent=!0},95:function(e,t,n){"use strict";n.d(t,"a",(function(){return p})),n.d(t,"b",(function(){return m}));var r=n(0),a=n.n(r);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function s(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?s(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):s(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function l(e,t){if(null==e)return{};var n,r,a=function(e,t){if(null==e)return{};var n,r,a={},o=Object.keys(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);for(r=0;r<o.length;r++)n=o[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var u=a.a.createContext({}),c=function(e){var t=a.a.useContext(u),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},p=function(e){var t=c(e.components);return a.a.createElement(u.Provider,{value:t},e.children)},d={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},b=a.a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,o=e.originalType,s=e.parentName,u=l(e,["components","mdxType","originalType","parentName"]),p=c(n),b=r,m=p["".concat(s,".").concat(b)]||p[b]||d[b]||o;return n?a.a.createElement(m,i(i({ref:t},u),{},{components:n})):a.a.createElement(m,i({ref:t},u))}));function m(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var o=n.length,s=new Array(o);s[0]=b;var i={};for(var l in t)hasOwnProperty.call(t,l)&&(i[l]=t[l]);i.originalType=e,i.mdxType="string"==typeof e?e:r,s[1]=i;for(var u=2;u<o;u++)s[u]=n[u];return a.a.createElement.apply(null,s)}return a.a.createElement.apply(null,n)}b.displayName="MDXCreateElement"}}]);