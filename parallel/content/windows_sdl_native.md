## SDL2 for native Windows

**If you already have Go and SDL2 installed in WSL2 and they are working, please ignore this guide.**

1. You should already have Go (native windows) https://golang.org/dl/ installed on your machine.

2. Download mingw-w64 [here](http://sourceforge.net/projects/mingw-w64/files/Toolchains%20targetting%20Win32/Personal%20Builds/mingw-builds/installer/mingw-w64-install.exe/download).

3. Open installer, Select `x86_64` (AMD64) as Architecture, and leave the rest as default.

4. The default destination folder should be `C:\Program Files\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0`, we will take this path as example.

5. Search `advanced system settings` , click on `Environment Variables`, in `System variables` list, double click on `Path`.

6. Add a new variable `C:\Program Files\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0\mingw64\bin`.

   (if you already have Haskell installed on your machine, you probably need to move the newly added variables to the top of the list using `Move Up` button)

7. To confirm, open a new **PowerShell** and type `(gcm gcc).path`, make sure the path is the same as above.

8. Download SDL2 https://www.libsdl.org/download-2.0.php, move on to `Development Libraries` `Windows`and **download `SDL2-devel-2.x.xx-mingw.tar.gz(MinGW 32/64-bit)`**.(make sure you are downloading **Development Libraries** for **MinGW** and not for Visual C++)

9. Unzip it and open `x86_64-w64-mingw32`,  copy these four folders `bin` `include` `lib` `share` to

   `C:\Program Files\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0\mingw64\`

   as well as 

   `C:\Program Files\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0\mingw64\x86_64-w64-mingw32` 

   (You should copy those four folders to **both** of the path above)

10. Open your coursework folder, start a terminal and type `go run .` 

    The system will download dependency packages according to go.sum and go.mod. This usually takes a few minutes for cgo to prepare for the first run. 

    If you see the followings showing up in the terminal, you are ready for the coursework.

    ```powershell
    Threads: 8
    Width: 512
    Height: 512
    Completed Turns 0       Quitting
    ```

    


