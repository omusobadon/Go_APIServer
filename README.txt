### 環境構築方法 ###
1.　クローンの生成
    git clone https://github.com/omusobadon/Go_APIServer.git

2.　/Go_APIServer内で以下のコマンドを実行してDBを同期（DB操作用のパッケージが生成される）
    go run github.com/steebchen/prisma-client-go db push


### NAT設定 ###
ip nat inside source static tcp 192.168.1.7 8080 interface GigabitEthernet8 8080

### /POST 形式 ###
{
    "cust_id": 100,
    "product_id": 2,
    "order_num": 43
}