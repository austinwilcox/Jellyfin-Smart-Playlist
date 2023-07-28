# Jellyfin Smart Playlist
I created this repo as an effort to keep track of recent songs that I have added to my home jellyfin server.

The inspiration came from Itunes when it was brand new. I was able to go into itunes and create smart playlists based on a lot of different criteria, like recently added.

Currently v0.1 of this application will support recently added tracks, and I have plans to add different types of smart playlists.

## How it works
The program executes a file watcher on your music directory, and recursively checks all files for changes (additions of files). When an addition is made, an entry is appended to the playlist file.

## Setup
To set this up (Currently on Linux and Mac are supported, due to how I setup the config location), create a config.txt file in ~/.config/jellyfin-smart-playlist/ and the following keys need to be supplied.
```
name_of_playlist={PATH_TO_PLAYLIST_DIRECTORY}/SmartPlaylist/playlist.xml
playlist_title=Smart Playlist v0.2
user_id={JELLYFIN_USER_ID}
can_edit=true
folder_to_watch={FOLDER_TO_WATCH}
sub_folder={SUB_FOLDER}
index_to_take_onwards={INDEX_TO_SWAP_FOLDER}
```
So my case might be different than most, I sync my home nextcloud with music, and that then is pulled into Jellyfin. So I need to use a subfolder in order to access everything.

The env variables that might get confusing are folder_to_watch, sub_folder, and index_to_take_onwards, So i'll explain with an example.
Lets say my folder to watch is: /home/john/Music, and in Jellyfin the route to music looks like this: /Jellyfin/Music then sub_folder will be /Jellyfin (going all the way up to the matching folder name), and index_to_take_onwards would be 3, grabbing the 3rd item, which is music, so that in the playlist file it is creating the correct file link.

After you have this all setup, the application can be run with 
```
go run .
```
Or you can compile it and execute it that way.
