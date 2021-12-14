ffmpeg \
-protocol_whitelist \
file,udp,rtp \
-i \
/Users/rizkyramadhan/Works/WEB/Services/livestream/sdpcollection/8112-rtpforwarder.sdp \
-c:v \
libx264 \
-g 120 \
-keyint_min 120 \
-b:v 4000k \
-vf "fps=30,drawtext=fontfile=utils/OpenSans-Bold.ttf:box=1:fontcolor=black:boxcolor=white:fontsize=100:x=40:y=500:textfile=text.txt" \
-format_options movflags=cmaf \
-write_prft 1 \
-sc_threshold 0 \
-method PUT \
-ldash 1 \
-seg_duration 4 \
-frag_duration 1 \
-streaming 1 \
-http_persistent 1 \
-utc_timing_url "https://time.akamai.com/?iso" \
-use_timeline 0 \
-use_template 1 \
-media_seg_name 'chunk-stream-$RepresentationID$-$Number%05d$.m4s' \
-init_seg_name 'init-stream1-$RepresentationID$.m4s' \
-window_size 5  \
-extra_window_size 10 \
-remove_at_exit 1 \
-adaptation_sets "id=0,streams=v id=1,streams=a" \
-fflags \
nobuffer \
-f \
dash \
https://bifrost.inlive.app/ldash/8112-5a2d1590dc014ab0a6ce344cecc2cccb/manifest.mpd