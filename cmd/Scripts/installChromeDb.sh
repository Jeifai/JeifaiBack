export CHROME_BIN=/usr/bin/google-chrome
export DISPLAY=:99.0
sh -e /etc/init.d/xvfb start
sudo apt-get update
sudo apt-get install -y libappindicator1 fonts-liberation libasound2 libgconf-2-4 libnspr4 libxss1 libnss3 xdg-utils
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo dpkg -i google-chrome*.deb
sudo apt --fix-broken install
