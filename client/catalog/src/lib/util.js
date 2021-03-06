export let cap= (s) => {
	if (!s) {
		return "";
	}
	s = s.replace(/([A-Z]+)/g, " $1").replace(/([A-Z][a-z])/g, " $1");
	return s.charAt(0).toUpperCase() + s.slice(1)
}

export let restrictToIdent= (s) => {
	const newVal = s.replace(/[^a-zA-Z_]/g, '');
	return newVal;
}

export let isObject = s => {
  return typeof s == "object";
};

export function isEmptyObject(o) {
  return JSON.stringify(o) === "{}";
};

export function nonEmpty(o) {
  return !isEmptyObject(o);
};

export function clone(o) {
	return JSON.parse(JSON.stringify(o));
}
