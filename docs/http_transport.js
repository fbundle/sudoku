function httpTransportJsName(method, path) {
    return method + "$" + path.replaceAll("/", "$");
}

export async function httpPost(path, data, cb) {
    const name = httpTransportJsName("POST", path);
    if (typeof window[name] === "function") {
        const result = JSON.parse(window[name](JSON.stringify(data)));
        if (cb) cb(result);
    } else {
        const resp = await fetch(path, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(data),
        });
        if (cb) cb(await resp.json());
    }
}
