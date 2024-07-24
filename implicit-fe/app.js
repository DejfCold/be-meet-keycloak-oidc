console.log("loaded")

function onLogin() {
    const params = new URLSearchParams()
    params.set("nonce", "1")
    params.set("response_type", "token id_token")
    params.set("scope", "openid profile")
    params.set("client_id", "be-meet-implicit")
    params.set("redirect_uri", `${location.origin}${location.pathname}`)
    location.href = `http://localhost:8080/realms/be-meet/protocol/openid-connect/auth?${params}`
}

function onLogout() {
    const params = new URLSearchParams()
    params.set("post_logout_redirect_uri", `${location.origin}${location.pathname}`)
    params.set("client_id", "be-meet-implicit")
    location.href = `http://localhost:8080/realms/be-meet/protocol/openid-connect/logout?${params}`
}

addEventListener("DOMContentLoaded", (event) => {
    const params = new URLSearchParams(location.hash.slice(1))
    console.log(params)
    const accessToken = params.get("access_token")
    const atParts = accessToken.split(".")
    const header = JSON.stringify(JSON.parse(atob(atParts[0])), null, "\t")
    const payload = JSON.stringify(JSON.parse(atob(atParts[1])), null, "\t")
    const signature = atParts[2]
    document.getElementById("ta").innerText = "Header:\n" +  header + "\n\nPayload:\n" + payload + "\n\n" + signature + "\n\n"
});
