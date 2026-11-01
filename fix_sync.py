import re
with open('D:/PortfolioProjects/helix/storage/bench_test.go', 'r') as f:
    content = f.read()
content = content.replace('sync"', '"sync"')
with open('D:/PortfolioProjects/helix/storage/bench_test.go', 'w') as f:
    f.write(content)
print('Done')
