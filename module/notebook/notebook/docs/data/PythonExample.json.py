# $PF_GENERATE_ONCE$
import json
import sys
from datetime import datetime

x = [["Python", datetime.now().isoformat()]]

json.dump(x, sys.stdout)
