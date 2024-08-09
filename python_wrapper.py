import sys
import subprocess
import requests
from datetime import datetime

def main():
    if len(sys.argv) < 3:
        print("Usage: python track_time.py <name> <command>")
        sys.exit(1)

    global name
    global command

    name = sys.argv[1]
    command = ' '.join(sys.argv[2:])

    # Record the start time
    start_time = datetime.now().strftime('%Y-%m-%dT%H:%M:%SZ')

    # POST start time
    requests.post('http://localhost:8080/track_start', data=f"{name},{start_time}")

    # Execute the command and capture its exit status
    try:
        subprocess.run(command, shell=True, check=True)
        status = 0
    except subprocess.CalledProcessError as e:
        status = e.returncode

    # Record the end time
    end_time = datetime.now().strftime('%Y-%m-%dT%H:%M:%SZ')

    # POST end time
    requests.post('http://localhost:8080/track_end', data=f"{name},{end_time}")

    # Exit with the status of the command
    sys.exit(status)

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        end_time = datetime.now().strftime('%Y-%m-%dT%H:%M:%SZ')

        requests.post('http://localhost:8080/track_end', data=f"{name},{end_time}")
