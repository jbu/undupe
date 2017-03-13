import exifread

photo_masters = '/Users/james/Pictures/Photos Library.photoslibrary/Masters/'
t = exifread.process_file(open(photo_masters + '2016/12/16/20161216-111800/DSC_0295.NEF','rb'))

for tag in t.items():
    print(tag[0], tag[1])
