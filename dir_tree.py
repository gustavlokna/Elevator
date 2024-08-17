import os

def generate_tree(root_dir, padding=""):
    print(padding + os.path.basename(root_dir) + "/")
    padding = padding + "    "
    items = sorted(os.listdir(root_dir))
    for index, item in enumerate(items):
        path = os.path.join(root_dir, item)
        if os.path.isdir(path):
            generate_tree(path, padding + "|  ")
        else:
            print(padding + "|-- " + item)

if __name__ == "__main__":
    generate_tree(".")
