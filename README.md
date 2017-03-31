http://www.pyimagesearch.com/2016/12/19/install-opencv-3-on-macos-with-homebrew-the-easy-way/
brew install opencv3 --with-contrib --with-python3 --HEAD
sudo mv /usr/local/opt/opencv3/lib/python3.6/site-packages/cv2.cpython-36m-darwin.so /usr/local/opt/opencv3/lib/python3.6/site-packages/cv2.so
echo /usr/local/opt/opencv3/lib/python3.6/site-packages >> /usr/local/lib/python3.6/site-packages/opencv3.pth
echo /usr/local/opt/opencv3/lib/python3.6/site-packages >> .e/lib/python3.6/site-packages/opencv3.pth
python3 -m venv .e
. .e/bin/activate
pip install -r requirements.txt
jupyter notebook unduplify.ipynb
