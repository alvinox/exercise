import runtime_path  # isort:skip

import argparse
import gzip
import os
import pickle
import time
from urllib.error import URLError
from urllib.request import urlretrieve

import numpy as np

from core.evaluator import AccEvaluator
from core.layers import Dense
from core.layers import ReLU
from core.losses import SoftmaxCrossEntropyLoss
from core.model import Model
from core.nn import Net
from core.optimizer import Adam
from utils.data_iterator import BatchIterator


def get_one_hot(targets, nb_classes):
    return np.eye(nb_classes)[np.array(targets).reshape(-1)]


def prepare_dataset(data_dir):
    if not os.path.exists(data_dir):
        os.makedirs(data_dir)

    url = "http://deeplearning.net/data/mnist/mnist.pkl.gz"
    path = os.path.join(data_dir, url.split("/")[-1])

    # download
    try:
        if os.path.exists(path):
            print("{} already exists.".format(path))
        else:
            print("Downloading {}.".format(url))
            try:
                urlretrieve(url, path)
            except URLError:
                raise RuntimeError("Error downloading resource!")
            finally:
                print()
    except KeyboardInterrupt:
        print("Interrupted")

    # load
    print("Loading MNIST dataset.")
    with gzip.open(path, "rb") as f:
        return pickle.load(f, encoding="latin1")


def main(args):
    train_set, valid_set, test_set = prepare_dataset(args.data_dir)
    train_x, train_y = train_set
    test_x, test_y = test_set
    train_y = get_one_hot(train_y, 10)

    net = Net([
        Dense(784, 200),
        ReLU(),
        Dense(200, 100),
        ReLU(),
        Dense(100, 70),
        ReLU(),
        Dense(70, 30),
        ReLU(),
        Dense(30, 10)
    ])

    model = Model(net=net, loss=SoftmaxCrossEntropyLoss(), optimizer=Adam(lr=args.lr))

    iterator = BatchIterator(batch_size=args.batch_size)
    evaluator = AccEvaluator()
    loss_list = list()
    for epoch in range(args.num_ep):
        t_start = time.time()
        for batch in iterator(train_x, train_y):
            pred = model.forward(batch.inputs)
            loss, grads = model.backward(pred, batch.targets)
            model.apply_grad(grads)
            loss_list.append(loss)
        t_end = time.time()
        # evaluate
        test_pred = model.forward(test_x)
        test_pred_idx = np.argmax(test_pred, axis=1)
        test_y_idx = np.asarray(test_y)
        res = evaluator.evaluate(test_pred_idx, test_y_idx)
        print("Epoch %d time cost: %.4f\t %s" % (epoch, t_end - t_start, res))


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--num_ep", default=50, type=int)
    parser.add_argument("--data_dir", default="./examples/mnist/data", type=str)
    parser.add_argument("--lr", default=1e-3, type=float)
    parser.add_argument("--batch_size", default=128, type=int)
    parser.add_argument("--seed", default=0, type=int)
    args = parser.parse_args()
    main(args)