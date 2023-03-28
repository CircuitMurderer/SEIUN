cp /etc/hosts /etc/hosts.bak

echo "127.0.0.1       council.seiun.net" >> /etc/hosts
echo "127.0.0.1       company.seiun.net" >> /etc/hosts
echo "127.0.0.1       school.seiun.net" >> /etc/hosts
echo "127.0.0.1       group.seiun.net" >> /etc/hosts

echo "127.0.0.1       orderer1.council.seiun.net" >> /etc/hosts
echo "127.0.0.1       orderer2.council.seiun.net" >> /etc/hosts
echo "127.0.0.1       orderer3.council.seiun.net" >> /etc/hosts

echo "127.0.0.1       peer1.company.seiun.net" >> /etc/hosts
echo "127.0.0.1       peer1.school.seiun.net" >> /etc/hosts
echo "127.0.0.1       peer1.group.seiun.net" >> /etc/hosts
