---
title: Ubuntu(Debian)でMQTTを使う方法
image: images/Microsoft-Fluentui-Emoji-3d-Factory-3d.1024.png
publishedAt: 2023-07-12T00:00:00.000Z
updatedAt: 2024-05-01T00:00:00.000Z
---
# MQTTとは

MQTTは、たくさんの端末同士で、短いメッセージを効率的に送受信するためのプロトコルです。複数のセンサーから送信されるデータを監視し、送信されたデータの値に応じて処理を実行できます。

# MQTTの歴史

1999年にIBMとEurotechが開発しました。当初は石油やガスの産業で使用することを想定しており、衛星経由で石油パイプラインをモニタリングするためには、帯域幅や消費電力を削減する必要があったため、それを達成するためのプロトコルとしてMQTTが生まれました。2010年には、IBMがMQTTをオープンプロトコルとして公開し、2013年にはOASISという標準化団体がMQTTの標準化を開始しました。

# MQTTの特徴

## データが軽量

MQTTはパケットのヘッダーサイズが小さいため、HTTPなどの他のプロトコルを使用した場合と比べてデータ量、CPUの負荷、消費電力が小さくなります。そのため、特にIoTアプリケーションと相性が良く、IoTの分野で普及しています。

## Pub/Sub(Publish-Subscribe)モデル

Pub/Subモデルは、たくさんの端末同士で、非同期にメッセージを送受信する仕組みです。送信されたメッセージは、トピックごとにブローカー(Broker)と呼ばれる中継システムに保管されます。メッセージの送信者はパブリッシャー(Publisher)と呼ばれ、メッセージの受信者はサブスクライバー(Subscriber)と呼ばれます。

### Pub/Subモデルの利点

Pub/Subモデルでは、サブスクライバーがトピックを後からサブスクライブしてメッセージを受信することもできるため、参加者の変化に強く、システムの規模の変更に対応しやすいという利点があります。

## Retained messagesが使える

Retained messagesは、それぞれのトピックに対して最後に送信されたメッセージを、ブローカーに保持する機能です。この機能を使用するには、RetainedフラグをTrueに設定します。

## QoSを設定できる

QoS(Quality of Service)は、通信の品質を制御する仕組みです。通信の品質が低すぎると適切に通信できず、高すぎると通信の無駄が多くなってしまいます。3種類の中から適切なQoSを設定することで通信の品質を制御します。

### QoSで制御される通信

パブリッシャーとブローカーの間の通信と、ブローカーとサブスクライバーの間の通信がQoSによって制御されます。

### QoS=0

受信側は、あるメッセージを最大で1回受信します。つまり、あるメッセージが受信されることもあれば、受信されないこともあります。通信エラーによってメッセージの受信に失敗した場合に、メッセージの再送は行われません。1番単純かつ軽量な通信方式です。

![](images/qos0.png)

引用: [https://sakiot.com/what-is-qos-of-mqtt/](https://sakiot.com/what-is-qos-of-mqtt/)

### QoS=1

受信側は、あるメッセージを最低1回受信します。つまり、あるメッセージが1回だけ受信されることもあれば、重複して受信されることもあります。QoS=1ではメッセージの重複を排除できないため、もし重複していたとしても異なるメッセージとして扱われます。

#### メッセージの再送で重複が発生しない場合

![](images/qos1-1.png)

引用: [https://sakiot.com/what-is-qos-of-mqtt/](https://sakiot.com/what-is-qos-of-mqtt/)

#### メッセージの再送で重複が発生する場合

![](images/qos1-2.png)

引用: [https://sakiot.com/what-is-qos-of-mqtt/](https://sakiot.com/what-is-qos-of-mqtt/)

### QoS=2

受信側は、あるメッセージを1回だけ受信します。通信の品質が最も高い通信方法ですが、通信回数が多いため比較的低速です。

#### 通信エラーが発生しない場合

![](images/qos2-1.png)

引用: [https://sakiot.com/what-is-qos-of-mqtt/](https://sakiot.com/what-is-qos-of-mqtt/)

#### 通信エラーが発生する場合

![](images/qos2-2.png)

引用: [https://sakiot.com/what-is-qos-of-mqtt/](https://sakiot.com/what-is-qos-of-mqtt/)

# MQTTを使った通信を試してみる

MQTTによる通信を行うためには、MQTTブローカーが必要です。主なMQTTブローカーとしては、Mosquitto, RabbitMQ, Apache ActiveMQがあります。今回はMosquittoをインストールします。

```bash
sudo apt install mosquitto mosquitto-clients
```

mosquitto.serviceが起動していることを確認します。起動していれば、activeと表示されます。

```bash
sudo systemctl is-active mosquitto
```

## helloトピックでメッセージを送受信

最初に、helloトピックに挨拶メッセージを送受信してみます。

まずはmosquitto\_subコマンドでトピックをサブスクライブします。

```bash
mosquitto_sub -h localhost -t hello
```

次に、mosquitto\_pubコマンドでメッセージをトピックにパブリッシュします。

```bash
mosquitto_pub -h localhost -t hello -m "Hello, world!"
```

## Retained messagesを使う

MosquittoでRetained messagesを使うには、-rオプションを指定します。

```
mosquitto_pub --help
-r: message should be retained.
```

```
mosquitto_pub -h localhost -t hello -m "This is a retained message" -r
```

Retained messagesはブローカーに保持されているので、それをパブリッシュした後でトピックをサブスクライブしてもメッセージを受信できます。

```
nagutabby@debian:~$ mosquitto_sub -h localhost -t hello
This is a retained message
```

## Retained messagesを削除する

Retained messagesを削除するには、空のRetained messagesを送信します。

```
mosquitto_pub --help
-n: send a null(zero length) message.
```

```
mosquitto_pub -h localhost -t hello -r -n
```

## トピックをまとめてサブスクライブする

#や+をサブスクライブ時に指定すると、トピックをまとめてサブスクライブできます。

まずはトピックとして、greeting/morningとgreeting/afternoonを作成します。

```
mosquitto_pub -h localhost -t greeting/morning -m "Good morning!" -r
mosquitto_pub -h localhost -t greeting/afternoon -m "Hello!" -r
```

次に、greetingという階層にあるトピックを#を使ってまとめてサブスクライブします。

```
mosquitto_sub -h localhost -t greeting/#
Good morning!
Hello!
```

## 気象データを送受信する

### OpenWeather

気象データを取得することができます。アカウント作成が必須です。

アカウント作成:

[https://home.openweathermap.org/users/sign\_up](https://home.openweathermap.org/users/sign_up)

APIキー:

[https://home.openweathermap.org/api\_keys](https://home.openweathermap.org/api_keys)

例えば、東京駅の気象データを取得するには、以下のコマンドを実行します。

```
curl "https://api.openweathermap.org/data/2.5/onecall?lat=35.681236&lon=139.767125&units=metric&lang=ja&appid={YOUR API KEY}"
```

### Open-Meteo

こちらでも気象データを取得できます。アカウント作成が不要であるため、今回はこちらを使用します。

ドキュメント: [https://open-meteo.com/en/docs](https://open-meteo.com/en/docs)

東京の1日の最高気温を直近1週間だけ取得してみます。

```
curl "https://api.open-meteo.com/v1/forecast?latitude=35.6895&longitude=139.6917&daily=temperature_2m_max&timezone=Asia%2FTokyo"
```

データを確認しやすくするために、レスポンスをjqコマンドで整形します。

```
curl "https://api.open-meteo.com/v1/forecast?latitude=35.6895&longitude=139.6917&daily=temperature_2m_max&timezone=Asia%2FTokyo" | jq
```

temperature\_2m\_maxの先頭の値を取り出し、パブリッシュするシェルスクリプトを書きます。

```
#!/bin/bash

temp=`curl "https://api.open-meteo.com/v1/forecast?latitude=35.6895&longitude=139.6917&daily=temperature_2m_max&timezone=Asia%2FTokyo" |  jq .daily  | jq .temperature_2m_max | jq ".[0]"`
echo $temp
mosquitto_pub -h localhost -t weather/tokyo/temperature_2m_max -m $temp -r
```

pub\_temp.shという名前で保存して実行します。そして、メッセージがパブリッシュされたトピックをサブスクライブします。

```
bash pub_temp.sh
mosquitto_sub -h localhost -t weather/tokyo/temperature_2m_max
```
