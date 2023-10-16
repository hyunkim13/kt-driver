curl -k -v -s -X POST "https://api.ucloudbiz.olleh.com/d1/identity/auth/tokens" \
--header "Content-Type: application/json" \
--data-raw '{"auth": {"identity": {"methods": ["password"], "password": { "user": { "domain": {"id": "default"},"name": "saasify@systeer.com", "password": "Sysmaster77!!"}}}, "scope": {"project": {"domain": {"id": "default"}, "name": "saasify@systeer.com"}}}}' | grep -i X-Subject-Token
