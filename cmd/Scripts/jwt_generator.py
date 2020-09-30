# You'll need to install PyJWT via pip 'pip install PyJWT' or your project packages file

import jwt

METABASE_SITE_URL = "http://metabase.jeifai.com"
METABASE_SECRET_KEY = "YOUR_SECRET_KEY"

payload = {
  "resource": {"dashboard": 2},
  "params": {
    
  }
}
token = jwt.encode(payload, METABASE_SECRET_KEY, algorithm="HS256")

iframeUrl = METABASE_SITE_URL + "/embed/dashboard/" + token.decode("utf8")

print iframeUrl