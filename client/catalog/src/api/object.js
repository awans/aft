import { mergeOps } from './merge.js'

function empty(o) {
	if (Array.isArray(o)) {
		return o.length === 0;
	} else if (typeof o === "object" && o !== null) {
		const ks = Object.keys(o);
		if (Array.isArray(ks) && (ks.length === 0 || ks.length === 1 && ks[0] === "type") ) {
			return true;
		}
	}
	return false;
};

function nonEmpty(o) {
	return !empty(o);
};

function clone(o) {
	return JSON.parse(JSON.stringify(o));
}


export function ObjectOperation(config) {
	const childOps = new Set();

	let init = null;

	const base = {
		__op: true,
		op: () => {
			const data = {};

			for (let key of childOps) {
				const childOp = base[key].op();
				if (childOp !== null && nonEmpty(childOp)) {
					data[key] = childOp;
				}
			}

			if (empty(data)) {
				return {};
			}

			let op;

			if (init) {
				op = {
					update: {
						where: {id: init.id},
						data: data,
					}
				}
			} else {
				op = {
					create: {
						data: data,
					}
				}
			}

			return op;
		},
		initialize: (i) => {
			init = JSON.parse(JSON.stringify(i));
			for (let [k, v] of Object.entries(i)) {
				if (childOps.has(k)) {
					base[k].initialize(v)
				}
			}
		},
		clientInit: (iVal) => {
			for (let [k, v] of Object.entries(iVal)) {
				if (childOps.has(k)) {
					base[k].clientInit(v)
				}
			}
		},
		clone: () => {
			const newConfig = {}
			for (let [k, v] of Object.entries(config)) {
				newConfig[k] = v.clone();
			}
			return ObjectOperation(newConfig);
		}
	};
	
	for (let [k,v] of Object.entries(config)) {
		if (v.__op) {
			childOps.add(k);
		} 
		base[k] = v;
	}


	return new Proxy(base, {
		set: function(target, prop, newVal) {
			if (childOps.has(prop)) {
				const d = target[prop]
				if (d.__descriptor) {
					d.set(newVal)
					return true;
				} else if (d.__op && newVal.__op) {
					target[prop] = newVal;
					return true;
				}
			}
		},
		get: function(target, prop) {
			const v = target[prop];
			if (v && v.__descriptor) {
				return v.get();
			}
			return v;
		}
	});
}

export function RelationshipOperation(config) {
	let base = []
	let proxy = new Proxy(base, {
		set: function(target, prop, newVal) {
			let p = target[prop]
			if (p && p.__descriptor) {
				p.set(newVal)
			} else {
				target[prop] = newVal;
			}
			return true;
		},
		get: function(target, prop) {
			const v = target[prop];
			if (v && v.__descriptor) {
				return v.get();
			}
			return v;
		}
	});

	base.__op = true
	base.init = null;
	base.op = () => {
		let ops = [];
		for (let child of base) {
			const childOp = child.op();
			if (nonEmpty(childOp)) {
				ops.push(childOp);
			}
		}
		return mergeOps(ops);
	}

	base.initialize = (ivals) => {
		base.init = JSON.parse(JSON.stringify(ivals));
		for (let v of ivals) {
			const c = config.clone();
			c.initialize(v);
			base.push(c);
		}
	}
	base.removeBy = (f) => {
		const ix = proxy.findIndex(f);
		if (ix !== -1){
			proxy[ix] = null;
		}
		return proxy
	}

	base.add = (clientInit) => {
		let n = config.clone()
		if (clientInit) {
			n.clientInit(clientInit);
		}
		base.push(n)
		return proxy
	}

	base.clone = () => {
		const newConfig = {}
		for (let [k, v] of Object.entries(config)) {
			newConfig[k] = v.clone();
		}
		return RelationshipOperation(newConfig);
	}
	return proxy
}

export function AttributeOperation(def) {
	const descriptor = {
		__op: true,
		__descriptor: true,
		init: null,
		value: def,
		get: function() {
			return descriptor.value;
		},
		set: function(newVal) {
			descriptor.value = newVal
		},
		initialize: function(iVal) {
			descriptor.init = iVal;
			descriptor.value = iVal;
		},
		clientInit: function(iVal) {
			descriptor.set(iVal);
		},
		op: function() {
			if (descriptor.init !== null && descriptor.value !== descriptor.init) {
				return descriptor.value;
			} else if (descriptor.init === null) {
				return descriptor.value;
			}

			return null
		},
		clone: function() {
			return AttributeOperation(def);
		}
	}
	return descriptor;
}


export function TypeSpecifier(ifaceName) {
	const descriptor = {
		__op: true,
		__descriptor: true,
		value: ifaceName,
		get: function() {
			return descriptor.value;
		},
		set: function(newVal) {
			return false;
		},
		initialize: function(iVal) {
			descriptor.value = iVal;
		},
		clientInit: function(iVal) {},
		op: function() {
			return descriptor.value;
		},
		clone: function() {
			return TypeSpecifier(ifaceName);
		}
	}
	return descriptor;
}

export function ConnectOperation(clientInit) {
	return RelOperation("connect", "disconnect", clientInit);
}

export function SetOperation(clientInit) {
	return RelOperation("set", "set", clientInit)
}

function RelOperation(opType, inverse, clientInit) {
	const descriptor = {
		__op: true,
		__descriptor: true,
		init: null,
		value: null,
		get: function() {
			return descriptor.value;
		},
		set: function(newVal) {
			descriptor.value = newVal;
		},
		initialize: function(iVal) {
			descriptor.init = clone(iVal);
			descriptor.value = iVal;
		},
		clientInit: function(iVal) {
			descriptor.value = iVal;
		},
		op: function() {
			if (descriptor.init && descriptor.value && 
				descriptor.value.id !== descriptor.init.id) {
				const op = {};
				op[opType] = {id: descriptor.value.id}
				return op;
			} else if (descriptor.init && descriptor.value === null) {
				const op = {};
				op[inverse] = {id: descriptor.init.id}
				return op;
			} else if (descriptor.init === null && descriptor.value) {
				const op = {};
				op[opType] = {id: descriptor.value.id}
				return op;
			}
			return null;
		},
		clone: function() {
			return RelOperation(opType, inverse);
		}
	}
	if (clientInit) {
		descriptor.clientInit(clientInit);
	}
	return descriptor;
}

export function OpLiteral(lit) {
	const descriptor = {
		__op: true,
		value: {},
		set: function(newVal) {
			return false;
		},
		op: function() {
			return lit
		},
		clone: function() {
			return OpLiteral(lit);
		}
	}
	return descriptor;
}

export function ReadOnly(defaultVal) {
	const descriptor = {
		__op: true,
		__descriptor: true,
		value: defaultVal,
		get: function() {
			return descriptor.value;
		},
		set: function(newVal) {
			return false;
		},
		initialize: function(iVal) {
			descriptor.value = iVal;
		},
		clientInit: function(iVal) {
			descriptor.value = iVal;
		},
		op: function() {
			return null
		},
		clone: function() {
			return ReadOnly(defaultVal);
		}
	}
	return descriptor;
}

export function Case(cases) {
	let type = null;

	const descriptor = {
		__descriptor: true,
		__op: true,
		value: null,
		get: function() {
			return descriptor.value;
		},
		set: function(prop, newVal) {
			return false;
		},
		clientInit: function(iVal) {
			type = iVal.type;
			descriptor.value = cases[type];
			descriptor.value.clientInit(iVal);
		},
		initialize: function(iVal) {
			type = iVal.type;
			descriptor.value = cases[type];
			descriptor.value.initialize(iVal);
		},
		op: function() {
			return descriptor.value.op();
		},
		clone: function() {
			const newCases = {};
			for (let [k, v] of Object.entries(cases)) {
				newCases[k] = v.clone();
			}
			return Case(newCases);
		}
	}
	return descriptor;
}