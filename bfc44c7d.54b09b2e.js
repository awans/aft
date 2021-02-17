(window.webpackJsonp=window.webpackJsonp||[]).push([[14],{84:function(e,t,n){"use strict";n.r(t),n.d(t,"frontMatter",(function(){return r})),n.d(t,"metadata",(function(){return i})),n.d(t,"toc",(function(){return l})),n.d(t,"default",(function(){return d}));var o=n(3),a=(n(0),n(95));const r={id:"models",title:"Models"},i={unversionedId:"tutorial/models",id:"tutorial/models",isDocsHomePage:!1,title:"Models",description:"Let's add another component in app.js to for our Todo data and reference it from our App component so it displays once the user is signed in. We'll also move our signout function into the Todos component, now that it has a better home.",source:"@site/docs/tutorial/06-models.md",slug:"/tutorial/models",permalink:"/docs/tutorial/models",version:"current",sidebar:"main",previous:{title:"User",permalink:"/docs/tutorial/user"},next:{title:"Creating",permalink:"/docs/tutorial/creates"}},l=[{value:"Adding Models",id:"adding-models",children:[]},{value:"Reading data",id:"reading-data",children:[]}],s={toc:l};function d({components:e,...t}){return Object(a.b)("wrapper",Object(o.a)({},s,t,{components:e,mdxType:"MDXLayout"}),Object(a.b)("p",null,"Let's add another component in ",Object(a.b)("inlineCode",{parentName:"p"},"app.js")," to for our Todo data and reference it from our App component so it displays once the user is signed in. We'll also move our signout function into the Todos component, now that it has a better home."),Object(a.b)("pre",null,Object(a.b)("code",Object(o.a)({parentName:"pre"},{className:"language-js",metastring:'title="app.js"',title:'"app.js"'}),'export function App (props) {\n    ...\n    if (!loaded) {\n        return html``\n    } else if (user === null) {\n        return html`<${Login} setUser=${setUser} />`\n    } else {\n        return html`<${Todos} user=${user} setUser=${setUser}/>`\n    }\n}\n\nfunction Todos({user, setUser}) {\n    const [todos, setTodos] = useState([]);\n\n    const signout = async () => {\n        await aft.rpc.logout();\n        setUser(null);\n    }\n\n    return html`\n    <div class="box">\n        <div class="row">\n            <div><b>Todos</b></div><a onClick=${signout}>Sign out</a>\n        </div>\n        ${todos.map(todo => {\n            return html`<${Todo} key=${todo.id} todo=${todo} />`\n        })}\n    </div>`\n}\n\nfunction Todo({todo}) {\n    return html`\n    <div class="row">\n        <div>${todo.text}</div>\n        <input type=checkbox checked=${todo.done} />\n    </div>`\n}\n\n')),Object(a.b)("h2",{id:"adding-models"},"Adding Models"),Object(a.b)("p",null,"Now let's create a basic Todo object on our backend. We'll have two attributes, ",Object(a.b)("inlineCode",{parentName:"p"},"text"),", a string which is the text of the Todo, and ",Object(a.b)("inlineCode",{parentName:"p"},"done")," a boolean indicating whether it's done or not. We'll also add a relationship to a ",Object(a.b)("inlineCode",{parentName:"p"},"user")," object\u2014the owner of the Todo."),Object(a.b)("p",null,"Switch back over to Aft, and navigate to the ",Object(a.b)("strong",{parentName:"p"},"Schema")," section."),Object(a.b)("p",null,"Click the ",Object(a.b)("strong",{parentName:"p"},"Add Model")," button at the top of the screen. "),Object(a.b)("p",null,"Fill in the text field that says ",Object(a.b)("em",{parentName:"p"},"Model name..")," with ",Object(a.b)("inlineCode",{parentName:"p"},"todo"),","),Object(a.b)("p",null,"Then click the ",Object(a.b)("strong",{parentName:"p"},"add")," button under Attributes and fill in our first attribute."),Object(a.b)("p",null,"For ",Object(a.b)("em",{parentName:"p"},"Attribute name..")," type ",Object(a.b)("inlineCode",{parentName:"p"},"text"),", and select ",Object(a.b)("inlineCode",{parentName:"p"},"String")," from the dropdown on the right, indicating its type."),Object(a.b)("p",null,"Click ",Object(a.b)("strong",{parentName:"p"},"add")," again, and this time name the attribute ",Object(a.b)("inlineCode",{parentName:"p"},"done"),", and select ",Object(a.b)("inlineCode",{parentName:"p"},"Bool")," from the type dropdown."),Object(a.b)("p",null,"Now let's add our relationship to ",Object(a.b)("inlineCode",{parentName:"p"},"user")," by clicking the ",Object(a.b)("strong",{parentName:"p"},"add")," button under Relationships."),Object(a.b)("p",null,"Set the ",Object(a.b)("em",{parentName:"p"},"Relationship name..")," to ",Object(a.b)("inlineCode",{parentName:"p"},"user")," and select ",Object(a.b)("inlineCode",{parentName:"p"},"User")," from the dropdown. That tells aft that this property points to a User object."),Object(a.b)("p",null,"You can leave the multiple box unchecked, since this relationship will only be to a single user, rather than a list of users."),Object(a.b)("p",null,"Once you've done that, click ",Object(a.b)("strong",{parentName:"p"},"Save")," next to the model name at the top of the page, and you're all done!"),Object(a.b)("h2",{id:"reading-data"},"Reading data"),Object(a.b)("p",null,"Aft automatically adds our new Todo model to the API, so lets try and read some data from it."),Object(a.b)("p",null,"First, navigate over to ",Object(a.b)("strong",{parentName:"p"},"Terminal")," in Aft and run the following function to add a Todo."),Object(a.b)("pre",null,Object(a.b)("code",Object(o.a)({parentName:"pre"},{className:"language-python"}),'def main():\n    return create("todo", {"data": {\n            "text":"connect the backend", \n            "done": False, \n            "user": {"connect": {"email": "user@example.com"}}\n        }})\n')),Object(a.b)("p",null,"And finally, let's update our Todos component to call the server!"),Object(a.b)("p",null,"Here we're making use the Aft API's ",Object(a.b)("inlineCode",{parentName:"p"},"findMany")," method. You'll note we don't filter on the currently signed in user\u2014we'll be doing that automatically later on when we look at access controls. "),Object(a.b)("pre",null,Object(a.b)("code",Object(o.a)({parentName:"pre"},{className:"language-js",metastring:'title="app.js"',title:'"app.js"'}),'...\n\nfunction Todos({user, setUser}) {\n    const [todos, setTodos] = useState([]);\n    const [loaded, setLoaded] = useState(false);\n\n    useEffect(async () => {\n        try {\n            setTodos(await aft.api.todo.findMany({\n                where: {\n                    done: false,\n                }\n            }));\n        } catch {\n        } finally {\n            setLoaded(true);\n        }\n    }, []);\n\n    const signout = async () => {\n        await aft.rpc.logout();\n        setUser(null);\n    }\n\n    if (!loaded) {\n        return html``\n    }\n\n    return html`\n    <div class="box">\n        <div class="row">\n            <div><b>Todos</b></div><a onClick=${signout}>Sign out</a>\n        </div>\n        ${todos.map(todo => {\n            return html`<${Todo} key=${todo.id} todo=${todo} />`\n        })}\n    </div>`\n}\n\n...\n')),Object(a.b)("p",null,"Hit refresh and you should see the Todo you created on the server!"),Object(a.b)("p",null,"In the next section, we'll look at how we can add some data from the client."))}d.isMDXComponent=!0},95:function(e,t,n){"use strict";n.d(t,"a",(function(){return p})),n.d(t,"b",(function(){return m}));var o=n(0),a=n.n(o);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function i(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var o=Object.getOwnPropertySymbols(e);t&&(o=o.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,o)}return n}function l(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?i(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):i(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function s(e,t){if(null==e)return{};var n,o,a=function(e,t){if(null==e)return{};var n,o,a={},r=Object.keys(e);for(o=0;o<r.length;o++)n=r[o],t.indexOf(n)>=0||(a[n]=e[n]);return a}(e,t);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);for(o=0;o<r.length;o++)n=r[o],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(a[n]=e[n])}return a}var d=a.a.createContext({}),c=function(e){var t=a.a.useContext(d),n=t;return e&&(n="function"==typeof e?e(t):l(l({},t),e)),n},p=function(e){var t=c(e.components);return a.a.createElement(d.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.a.createElement(a.a.Fragment,{},t)}},b=a.a.forwardRef((function(e,t){var n=e.components,o=e.mdxType,r=e.originalType,i=e.parentName,d=s(e,["components","mdxType","originalType","parentName"]),p=c(n),b=o,m=p["".concat(i,".").concat(b)]||p[b]||u[b]||r;return n?a.a.createElement(m,l(l({ref:t},d),{},{components:n})):a.a.createElement(m,l({ref:t},d))}));function m(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var r=n.length,i=new Array(r);i[0]=b;var l={};for(var s in t)hasOwnProperty.call(t,s)&&(l[s]=t[s]);l.originalType=e,l.mdxType="string"==typeof e?e:o,i[1]=l;for(var d=2;d<r;d++)i[d]=n[d];return a.a.createElement.apply(null,i)}return a.a.createElement.apply(null,n)}b.displayName="MDXCreateElement"}}]);