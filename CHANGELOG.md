### 10/14/20
* program can loop through hardcoded directory and determine if patient sub-directories need image recovery
* program prints user/env info at runtime

### 10/15/20 v0.1
* program only moves image files from sub-directory containing "Original"
* program cleans up empty sub-directories it comes accross
* efficiency increased
    - program no longer enters and exits sub-directory on a loop
    - sub-directory path "*OriginalImages.XVA*" now referenced as global variable

### 10/15/20 v0.2
* program accepts -path flag
    - path flag can be used to run program in specified filepath
    - example: *./image_recovery -path="Z:\imaging-software\share\patientimages"*
* performance increase
    - *string.ReplaceAll* used instead of *Strings.Replace*