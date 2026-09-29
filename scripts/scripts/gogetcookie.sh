#!/bin/bash

eval 'set +o history' 2>/dev/null || setopt HIST_IGNORE_SPACE 2>/dev/null
 touch ~/.gitcookies
 chmod 0600 ~/.gitcookies

 git config --global http.cookiefile ~/.gitcookies

 tr , \\t <<\__END__ >>~/.gitcookies
go.googlesource.com,FALSE,/,TRUE,2147483647,o,git-marinssalinas.gmail.com=1/fzxxUwt0Hs8AFhJUiZ-IzI2fe_w6azG7h_D7X3J3kH4
go-review.googlesource.com,FALSE,/,TRUE,2147483647,o,git-marinssalinas.gmail.com=1/fzxxUwt0Hs8AFhJUiZ-IzI2fe_w6azG7h_D7X3J3kH4
__END__
eval 'set -o history' 2>/dev/null || unsetopt HIST_IGNORE_SPACE 2>/dev/null
