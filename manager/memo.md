# manager  
* clientとexecutorのswitch的な役割

コンテナ/ジョブの管理を行う
  - container
    - build
    - remove 
    - status

  - job
    - enqueue
    - dequeue
    - allocate
    - status

## 流れ
* jobはキューに貯めていく
* デキューする(apiを作って置かないと、clientがredisのクライアントの機能を持たないといけないので面倒。)
* プログラムが書かれたファイルを作る。(コンテナIDと同名)
* イメージを作る。(docker build) idをもらう。
* コンテナを立ち上げる(docker build)
* コンテナ、ジョブのステータス管理はmysqlで(id, status)
* executorに対してコンテナ情報とjob情報を与える(gorutine)
* executorからjobの結果をもらい、clientに返す
* DBからコンテナ情報を削除


## todo
* dockerのクライアントライブラリ(biuld, run, rm, rmi)
* postでプログラムを受け付けるサーバー
* ステータス管理, jobの管理をするmanager
* jobの実行をするexecutor
* botの作成
