import asyncio
import json
import machine
import re

from configs import CONFIG_PATH
from measurement import Measurement  # your I2C device reading module
from utils import get_version

# ----------------------------
# Globals
# ----------------------------
t = None  # will store latest I2C reading
h = None  # will store latest I2C reading

# ----------------------------
# Config handling
# ----------------------------
def load_config():
    with open(CONFIG_PATH, "r", encoding="utf-8") as f:
        return json.load(f)

def save_config(cfg):
    with open(CONFIG_PATH, "w", encoding="utf-8") as f:
        json.dump(cfg, f)

def try_cast(val):
    if val.isdigit():
        return int(val)
    try:
        return float(val)
    except:
        return val

def update_dict(d, key, value):
    keys = key.split(".")
    for k in keys[:-1]:
        d = d[k]
    d[keys[-1]] = try_cast(value)

# ----------------------------
# HTML generation
# ----------------------------
def dict_to_form(d, prefix=""):
    html = ""
    for k, v in d.items():
        full_key = f"{prefix}{k}" if prefix == "" else f"{prefix}.{k}"
        if isinstance(v, dict):
            html += f"<fieldset><legend>{k}</legend>"
            html += dict_to_form(v, prefix=full_key)
            html += "</fieldset>"
        else:
            safe_val = str(v).replace('"', "&quot;")
            html += f"""
            <label>{full_key}:
                <input type="text" name="{full_key}" value="{safe_val}">
            </label><br>
            """
    return html

def html_home():
    return f"""
    <html><body>
        <h1>Device Home</h1>

        <h3>STH3X</h3>
        <p>temperature: {t if t is not None else 'No data yet'}</p>
        <p>humidity: {h if h is not None else 'No data yet'}</p>
        <form action="/config"><button type="submit">Edit Config</button></form>
        <form action="/reset"><button type="submit">Reset</button></form>
    </body></html>
    """

def html_config(data):
    return f"""
    <html><body>
        <h1>Edit Configuration</h1>
        <form method="POST" action="/config">
            {dict_to_form(data)}
            <button type="submit">Update</button>
        </form>
        <form action="/"><button type="submit">Back to Home</button></form>
    </body></html>
    """

def html_error(msg):
    return f"""
    <html><body>
        <h1>Error</h1>
        <p>{msg}</p>
        <form action="/"><button type="submit">Back to Home</button></form>
    </body></html>
    """

def html_version():
    return f"""<!DOCTYPE html>
<html>
    <head> <title>envoops</title> </head>
    <body>
        <p>{get_version()}</p>
    </body>
</html>
"""

# ----------------------------
# HTTP helpers
# ----------------------------
def urldecode(s):
    s = s.replace("+", " ")
    parts = s.split("%")
    res = parts[0]
    for part in parts[1:]:
        if len(part) >= 2:
            try:
                res += chr(int(part[:2], 16)) + part[2:]
            except:
                res += "%" + part
        else:
            res += "%" + part
    return res

def parse_post_body(body):
    result = {}
    for pair in body.split("&"):
        if "=" in pair:
            key, value = pair.split("=", 1)
            result[urldecode(key)] = urldecode(value)
    return result

async def send_html(writer, html, status="200 OK"):
    writer.write(f"HTTP/1.0 {status}\r\nContent-Type: text/html\r\n\r\n")
    writer.write(html)
    await writer.drain()

async def redirect(writer, location):
    writer.write(f"HTTP/1.0 303 See Other\r\nLocation: {location}\r\n\r\n")
    await writer.drain()

# ----------------------------
# Request handler
# ----------------------------
async def serve_client(reader, writer):
    led = machine.Pin(2, machine.Pin.OUT)  # Adjust to your LED pin
    try:
        led.on()
        request_line = await reader.readline()
        if not request_line:
            await writer.aclose()
            return

        method, path, _ = request_line.decode().split()

        # Read headers
        content_length = 0
        while True:
            header = await reader.readline()
            if header == b"\r\n":
                break
            m = re.match("Content-Length: (\d+)", header.decode(), re.IGNORECASE)
            if m:
                content_length = int(m.group(1))

        data = load_config()

        if method == "GET":
            if path == "/":
                await send_html(writer, html_home())

            elif path == "/config":
                await send_html(writer, html_config(data))

            elif path == "/version":
                await send_html(writer, html_version())

            elif path.startswith("/error"):
                msg = "Unknown error"
                if "?" in path:
                    msg = urldecode(path.split("?", 1)[1])
                await send_html(writer, html_error(msg))

            elif path == "/reset":
                await send_html(writer, "<html><body><h1>Resetting...</h1></body></html>")
                await writer.aclose()
                machine.reset()

            else:
                await send_html(writer, "<html><body>404 Not Found</body></html>", "404 Not Found")

        elif method == "POST" and path == "/config":
            body = (await reader.read(content_length)).decode()
            params = parse_post_body(body)
            try:
                for k, v in params.items():
                    update_dict(data, k, v)
                save_config(data)
                await redirect(writer, "/")
            except Exception as e:
                await redirect(writer, f"/error?{urldecode(str(e))}")

        else:
            await send_html(writer, "<html><body>405 Method Not Allowed</body></html>", "405 Method Not Allowed")

    except Exception as e:
        print("Error:", e)
    finally:
        await writer.aclose()
        led.off()

# ----------------------------
# Background task
# ----------------------------
async def periodic_measurement(interval_sec=10):
    global t
    global h
    measurement = Measurement()
    while True:
        try:
            t, h = measurement.run()
        except Exception as e:
            #f"Error: {e}"
            t = "-250"
            h = "-250"
        await asyncio.sleep(interval_sec)

# ----------------------------
# Main loop
# ----------------------------
async def main():
    asyncio.create_task(periodic_measurement(10))  # every 10 seconds
    server = await asyncio.start_server(serve_client, "0.0.0.0", 80)
    print("Server running on port 80")
    await server.serve_forever()
