## バックエンド

ゲーム側に通信の責任を持たせないよう、バックエンドサーバに情報を集めてからゲーム側に情報を送っている

### SenderとHubとReceiver

ウェブクライアントからコメントを送るSender、情報を集計、監視するHub、情報を受け取るゲームクライアントのReceiverという3つのモデルを持たせて開発を行った。

Senderは視聴者の数いるので、バックエンドサーバーにSenderの数WebSocketのコネクションが行われる。そのため、各コネクションに対してインスタンスを生成し、Goの得意分野である並列処理、ゴルーチンを用いて各コネクションを管理している。

Hubは複数のSenderや単一のReceiverに関する情報を監視したり、メッセージをReceiverに送ったりすることを行っている。

Receiverはゲームクライアントとの状態を管理している。Hubから送られてきたメッセージをソケットを通してゲームクライアントに送信している。Receiverとのコネクションを管理している。

![Sender-Hub-Receiver](./docs/assets/StshootWebSocket.drawio.svg)
### Room

複数のゲーム部屋に対応するための実装を行った。

部屋作成エンドポイントにアクセスされたとき、Hubのインスタンスが生成され、生成されたHubはHubManagerによって管理される。

![Room](./docs/assets/StshotWebSocketRoom.svg)