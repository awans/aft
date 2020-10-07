const apiBase = "/api"
const apiMethods = ["create", "update", "findOne", "findMany"];
const objects = ["gift"];

const viewBase = "/views";
const viewMethods = ["login", "signup"]

const basePath = "https://localhost:8080";

function getToken() {
	let cookie = document.cookie;
	if (cookie) {
		let tok = cookie.split('; ')
		.find(row => row.startsWith('tok'))
		.split('=')[1];
		return tok;
	}
	return "";
}


async function post(url, params) {
	if(typeof params === 'undefined')  {
		params = {};
	}
	try {
		const res = await fetch(url, {
			method: "POST",
			body: JSON.stringify(params),
			headers: new Headers({
				'Authorization': getToken(),
			}),
			credentials: 'include',
		});
		const responseBody = await res.json();
		if ("code" in responseBody) {
			return Promise.reject(responseBody);
		}
		if("data" in responseBody) {
			return responseBody.data;
		}
	} catch (err) {
		console.log(err);
		throw err;
	}

}

function api(objects, methods) {
	const a = {};
	for (let o of objects) {
		a[o] = {};
		for (let m of methods) {
			a[o][m] = (params) => {
				return post(basePath + apiBase + '/' + o + "." + m, params);
			}
		}
	}
	return a;
}

function views(methods) {
	const v = {};
	for (let m of methods) {
		v[m] = (params) => {
			return post(basePath + viewBase + '/' + m, params);
		}
	}
	return v;
}

const client = {
	api: api(objects, apiMethods),
	views: views(viewMethods),
}

export default client;
