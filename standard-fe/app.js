console.log("loaded")

function onLogin() {
    const params = new URLSearchParams()
    params.set("nonce", "1")
    params.set("response_type", "code")
    params.set("scope", "openid profile")
    params.set("client_id", "be-meet-standard")
    params.set("redirect_uri", `${location.origin}/code-handler.html`)
    location.href = `http://localhost:8080/realms/be-meet/protocol/openid-connect/auth?${params}`
}

function onLogout() {
    const params = new URLSearchParams()
    params.set("post_logout_redirect_uri", `${location.origin}`)
    params.set("client_id", "be-meet-standard")
    location.href = `http://localhost:8080/realms/be-meet/protocol/openid-connect/logout?${params}`
}

addEventListener("DOMContentLoaded", async (event) => {
    const params = new URLSearchParams(location.search)
    console.log(params)
    const code = params.get("code")
    document.getElementById("ta").innerText = "Code:\n" + code

    setTimeout(async () => {
        await fetchToken(code)
    }, 1000)
});

async function fetchToken(code) {
    const resp = await fetch("http://localhost:8080/realms/be-meet/protocol/openid-connect/token", {
        method: "POST",
        headers: {
          "content-type": "application/x-www-form-urlencoded"
        },
        body: new URLSearchParams({
            "grant_type": "authorization_code",
            "client_id": "be-meet-standard",
            "code": code,
            "redirect_uri": `${location.origin}/code-handler.html`
        })
    })

    const respBody = await resp.json()
    console.log(respBody)

    const accessToken = respBody["access_token"]
    const atParts = accessToken.split(".")
    const header = JSON.stringify(JSON.parse(atob(atParts[0])), null, "\t")
    const payload = JSON.stringify(JSON.parse(atob(atParts[1])), null, "\t")
    const signature = atParts[2]
    document.getElementById("ta2").innerText = "Header:\n" +  header + "\n\nPayload:\n" + payload + "\n\n" + signature + "\n\n"
}