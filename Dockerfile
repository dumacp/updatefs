FROM golang:1.13

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN apt-get update
RUN apt-get -y install netcat openssh-server acl
RUN useradd -m sftpuser
RUN mkdir /home/sftpuser/.ssh
RUN echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDYl+tz9UPuQf7lu+tt4+F2ndrUbu4X2FlI0WmdEaDYgGA+G76nRN19cJK/NIXWH3DJWJG6lmZq1GwBLRrRWvKvqeiLN2Gn0df+V+POoMSYaokoLwYU4OHUOfavd7kY7AvPmEV/as3SG+ApiEZy7paYDUID98+HTcJ02j7Jjw/lecDutUySZkrIhqGVeT2F+sPgtnJvAEYJRmU/UqqOwSYbL3nFn3zH+d+yqpKuQpe/3GwWFIgzS7+gaHH3MTmTY3DG5OWRJLTvewKUSYhBJDGxaBH6+g9DjG0ZQpirQ2Y5uxmKcBpcMNd6f5GJE1aIhYk1E0p/ZpQsdNmOCky11lPB laboratorio@laboratorio-NE" >> /home/sftpuser/.ssh/authorized_keys
RUN echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCr01QfxA4oFHp1aJXyg8v9qwuZA6eUNuzifJj84GER4DsiwJJRV8Y0klMHHAIGovmgTkV0lkQzs6QE3AJuRDk1EiI8F1cDy+kB5QHpmHE+LZphEQMp/+aIu1SK1dBwJvg5J96Ys1/It1B9YRZMdJ+suu3qa2T5yGdXhR2McIIhiG2DPebX7hRekxf3avd0pI8buMQxElNFArkOEh1DBkn1qqbm3ZKoc5jLCWwGPA6C/s7Asmhb41Mu1//jOHNehCi5ST5P4Wjb4Bu2p92obV6r1HwK4G/TLXKnFP7blfRzAYIrBPn3TW3D0wBczb+Ao7StmIOecqzeM9yr6j5vwpYt duma@localhost.localdomain" >> /home/sftpuser/.ssh/authorized_keys
RUN echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC1KGk2N6Dy+U8vESKDrE5z5ZlFDbOGUG+l21+leCNyUZseuxeKVi70Wvn2zC4SU5a/qQIcy1zSmRgGyCsRBFuH5hw8UzDsgedsXXLdz2c2V/574w7ihjZBVEf18rOT9dSM/cEXrqLFfO6Z0lMA9qt0y2T/zvuGq9plU4fZdz9Xn0BYrEIoAfV0Py+W9EvUH4TXvvHzFzTgDnsQHZBBVocoPp35R84/S8uHUCvtXA+C7csaBYFeDAySMDVcMb9DAK1Iy5VJ4h2chWg9126wy1Q+wuW4klfQFoKN9kgjKkCW8JArlxTLCc15zOFwiS7RhrmT8jlSIRkid9wOIrnEW/YrlvCTcjxkIc0wfQhA5t6gmylRbsvtBzesgq6XpWgQXUD4nRQGC1nNFKD5cFpds/ygngrgEvWfytqHUN0/ArrhP1DvgEXjB0HPg49ZVPW5Y61gbJ/cuWbljoPWsiOJnJtyWyf8LBwOeRwBE3cfQdd+UkwQp2xguy0/LA/5csbe9is= duma@pop-os" >> /home/sftpuser/.ssh/authorized_keys
RUN chown -R sftpuser.sftpuser /home/sftpuser/.ssh
RUN chmod 640 /home/sftpuser/.ssh/authorized_keys

VOLUME /data/all

EXPOSE 8000/tcp

ENTRYPOINT ["./entrypoint.sh"]
CMD ["updatefs"]
