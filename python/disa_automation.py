#!/bin/python3

import re, subprocess, sys, os, requests, typing, jinja2, xml, ssl, http

def find_pull_point(line: str = None):
    resutl = re.search(r'(?<=<{opening_tag}>).+(?=<closing_tag>)', line)
    if result:
        return result
    else:
        return False

def get_pull_points(script: str = None):
    result = subprocess.popen(f'{script}')
    pull_points = list()

    for line in result.readlines(): # Idk, whatever result you get from subrpocess to read stdout
        #pull_points.append(re.search(r'(?<=<opening_tag>).{32}(?=</closing_tag>\s*$)', line))
        result = find_pull_point(line)
        if result:
            pull_points.append(result)
    return pull_points

def  template_xml(xml_template_path = None, xml_pull_point = None, xml_topic = None):
    j2_env = jinja2.Environment(loader = jinja2.FileSystemLoader(os.path.dirname(xml_template_path)))
    j2_template = j2_env.get_template(xml_template_path.split('/')[-1])
    content = j2.render(topic = xml_topic, pull_point = xml_pull_point)
     with open(f'{xml_template_path}.txt', 'w') as xml_template:
        xml_template.write(content) 

def query_disa(cert = None, key = None, xml_file = None):
#    cert = cert_file
#    key = key_file
    context = ssl.SSLContext(ssl.PROTOCOL_TLSv12)
    context.load_cert_chain(certfile = cert, keyfile = key)
    request_headers = { 'Content-Type': 'text/xml'}
    connection = http.client.HTTPSConnection(host = host, port = 443, context = context)
    with open(xml_file, 'r') as xmlfile:
        response = connection.request(method="POST", url = url, headers = request_headers, body = xml.loads(xmlfile.read()))
        return response

def validate_disa_response(response: xml = None):
    for line in response.content:
        if re.search(r'accessdenied|serverfault', line, re.IGNORECASE):
            return False
    return True


 




'''
https://www.techcoil.com/blog/how-to-send-a-http-request-with-client-certificate-private-key-password-secret-in-python-3/
https://realpython.com/primer-on-jinja-templating/#get-started-with-jinja
https://stackoverflow.com/questions/70036236/how-do-i-get-the-pem-from-jks-file
https://www.ibm.com/docs/en/slac/10.2.0?topic=uxws-convert-user-keys-certificates-pem-format-python-clients
https://security.stackexchange.com/questions/226747/what-is-the-difference-between-a-certificate-and-a-private-key
https://www.misterpki.com/python-requests-authentication/
'''





        
