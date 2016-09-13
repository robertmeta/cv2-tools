# cv2-goparser
This is a go decoder/encoder for the cv2 resume format. It should be able
to take a cv2 formatted document and turn it into a go struct, and take
a properly formatted go struct and turn it into a .cv2 file.

Currently it uses map of strings to interface{} in the hackParser

Later it will provide a console app for doing various cv2 related things.

# Try it

    go get -u -v github.com/cv2me/cv2-tools/cv2-goparser && cv2-goparser

Until PR is accepted, replace cv2me above with temblortenor.

# Contribute

Make a fork of the project on GitHub first (post PR)

    go get github.com/cv2me/cv2-tools/cv2-goparser
    cd $GOPATH/src/github.com/cv2me/cv2-tools/cv2-goparser
    git remote add fork git@github.com/{YOU}/{thefork.git}
    # do stuff, then try it
    go install -v github.com/cv2me/cv2-tools/cv2-goparser && cv2-goparser
    # to share changes
    git push fork master

Then submit a PR from your project to mine.

# TODO
- fix \{ ... } handling with multiple in section (push, pop?)
- 2nd pass processing to connect like \email and stuff like that

# Notes
- has to be two pass in order to work (cause of binding things found later)

# Useful Links
- https://www.w3.org/community/cv2/2015/09/29/introduction-to-svg-templates/
- https://www.w3.org/community/cv2/2015/09/30/proposal-cv2-syntax/
- https://blog.gopheracademy.com/advent-2014/parsers-lexers/

# Questions
- Only need \{...} support on value side?
- Maybe a diff type of \[...] for escaping tags?
- How do I know what section to get like \email from?
