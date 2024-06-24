rm *.pem

# 1. Generate CA's private key and self-signed certificate
# NOTICE: -nodes for ignore passphrase asking for development only, remove it for production to enter your password to improve the security
openssl req -x509 -newkey rsa:4096 -days 365 -nodes -keyout ca-key.pem -out ca-cert.pem -subj "/C=VN/ST=Hanoi/L=Hanoi/O=hwebz/OU=hwebz/CN=hwebz/emailAddress=hadm.haui@gmail.com"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

# 2. Generate web server's private key and certificate signing request (CSR)
# NOTICE: -nodes for ignore passphrase asking for development only, remove it for production to enter your password to improve the security
openssl req -newkey rsa:4096 -nodes -keyout server-key.pem -out server-req.pem -subj "/C=VN/ST=Hochiminh/L=Hochiminh/O=hdevz/OU=hdevz/CN=hdevz/emailAddress=hdevz.com@gmail.com"

# 3. Use CA's private key to sign web server's CSR and get back the signed certificate
# NOTICE: for localhost, we have to put IP=0.0.0.0 in server-ext.cnf, for production, put the real server's IP address
# the lack of IP, client can't connect to server
openssl x509 -req -in server-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem -extfile server-ext.cnf

echo "Server's signed certificate"
openssl x509 -in server-cert.pem -noout -text

# 4. Generate client's private key and certificate signing request (CSR)
openssl req -newkey rsa:4096 -nodes -keyout client-key.pem -out client-req.pem -subj "/C=VN/ST=Haiphong/L=Haiphong/O=hawebz/OU=hawebz/CN=hawebz/emailAddress=hawebz.com@gmail.com"

# 5. Use CA's private key to sign client's CSR and get back the signed certificate
openssl x509 -req -in client-req.pem -days 60 -CA ca-cert.pem -CAkey ca-key.pem -CAcreateserial -out client-cert.pem -extfile client-ext.cnf

echo "Client's signed certificate"
openssl x509 -in client-cert.pem -noout -text