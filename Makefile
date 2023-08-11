# #############################################################################
# # File: Makefile                                                             #
# # Project: goframe                                                           #
# # Created Date: 2023/08/12 00:57:46                                          #
# # Author: realjf                                                             #
# # -----                                                                      #
# # Last Modified: 2023/08/12 00:57:48                                         #
# # Modified By: realjf                                                        #
# # -----                                                                      #
# # Copyright (c) 2023                                                         #
# #############################################################################


.PHONY: build
build:
	go build github.com/realjf/goframe



.PHONY: run
run:
	go run goframe.go
