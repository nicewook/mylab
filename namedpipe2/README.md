# namedpipe close case

Namedpipe에서 쓰는 쪽에서 계속 열고 닫으면 어떻게 될까? 

쓰는 쪽에서 파일을 Close만 해도 읽는 쪽에서 0을 읽는 것을 예상했었는데 그러지는 않았다. 
아래 참고 링크에서는 그렇게 말했는데 흠흠. OS마다 다르다고는 한다 

참고링크: https://stackoverflow.com/questions/15055065/o-rdwr-on-named-pipes-with-poll/17384067#17384067
논블락 참고링크: https://medium.com/@cpuguy83/non-blocking-i-o-in-go-bc4651e3ac8d
