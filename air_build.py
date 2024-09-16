import platform
import subprocess

if __name__ == "__main__":
    if platform.system() == "Windows":
        subprocess.run("go build -o ./tmp/main.exe .", shell=True)
    else:
        subprocess.run("go build -o ./tmp/main .", shell=True)