# apimanager
APIManager CLI Tool

go install github.com/skckadiyala/apimanager

apimanager login

apimanager create org -n Marvel -ed
apimanager create user -n 'AntMan' -l antman -o Marvel -r user -p antman 
apimanager create user -n 'Thor' -l thor -o Marvel -r admin
apimanager create user -n 'Captain America' -l captain -o Marvel -r admin
apimanager create user -n 'Iron Man' -l ironman -o Marvel -r oadmin
apimanager create app -n Avengers -o Marvel
apimanager create key -a Avengers
apimanager create oauth -a Avengers -c resources/cert.pem 

apimanager create api -n 'Captain America' -o 'Marvel' -f resources/swagger.json 
apimanager create proxy -n 'The First Avenger' -b 'Captain America' -c resources/cert.pem -o Marvel -s passthrough
apimanager create proxy -n 'The Winter Soldier' -b 'Captain America' -c resources/cert.pem -o Marvel -s apikey -a Avengers
apimanager create proxy -n 'Civil War' -b 'Captain America' -c resources/cert.pem -o Marvel -s oauth -a Avengers

apimanager create api -n 'Iron Man' -o 'Marvel' -f resources/swagger.json 
apimanager create proxy -n 'Iron Man' -b 'Iron Man' -c resources/cert.pem -o Marvel (default security passthrough) 
apimanager create proxy -n 'Iron Man 2' -b 'Iron Man' -c resources/cert.pem -o Marvel -s apikey -a Avengers
apimanager create proxy -n 'Iron Man 3' -b 'Iron Man' -c resources/cert.pem -o Marvel -s oauth -a Avengers

apimanager list orgs
apimanager list users
apimanager list apps
apimanager list keys -a 'Avengers'
apimanager list oauths -a 'Avengers'
apimanager list apis
apimanager list proxies

apimanager unpublish proxy -n 'The First Avenger'
apimanager unpublish proxy -n 'The Winter Soldier'
apimanager unpublish proxy -n 'Civil War'


apimanager delete proxy -n 'The First Avenger'
apimanager delete api -n 'Captain America'
apimanager delete key -k 90b349a0-482c-44d4-b79a-fec70c71f809 -a Asgard
apimanager delete oauth -k 6126028c-5230-4b3b-b27b-f28f26dae28b -a Asgard
apimanager delete app -n Asgard
apimanager delete user -n 'AntMan'
apimanager delete user -n Thor
apimanager delete user -n 'Iron Man'
apimanager delete org -n Marvel

