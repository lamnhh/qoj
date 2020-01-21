import os
import random as rnd


# Configurations
TEST_COUNT = 100
PROBLEM_NAME = "MED"


def create_dir(test_id):
  """
  Create directory for test with a specified ID

  Parameters:
    test_id(int): id of test to be generated

  Return:
    None
  """
  dir_name = "Test%02d" % test_id
  if not os.path.exists(dir_name):
    os.mkdir(dir_name)


def generate(test_id, input_file):
  """
  Generate test with specified id

  Parameters:
    test_id(int): id of test to be generated
    input_file(file): file instance to be written on

  Return:
    None
  """
  if test_id <= TEST_COUNT * 0.25:
    # Subtask 1
    n = rnd.randint(1, 5)
  elif test_id <= TEST_COUNT * 0.6:
    # Subtask 2
    n = rnd.randint(6, 10)
  else:
    # Subtask 3
    n = rnd.randint(1, int(2e5))

  s = [chr(97 + rnd.randint(0, 25)) for i in range(n)]
  t = [chr(97 + rnd.randint(0, 25)) for i in range(n)]
  while (ord(s[n - 1]) + ord(t[n - 1])) % 2 == 1:
    t[n - 1] = chr(97 + rnd.randint(0, 25))

  s = "".join(s)
  t = "".join(t)
  if s > t:
    [s, t] = [t, s]
  
  input_file.write("%d\n" % n)
  input_file.write("".join(s) + "\n")
  input_file.write("".join(t) + "\n")


if __name__ == "__main__":
  for test_id in range(TEST_COUNT):
    # Create directory
    create_dir(test_id)

    # Generate inputs
    input_file = open("Test%02d/%s.INP" % (test_id, PROBLEM_NAME), "w+")
    generate(test_id, input_file)
    input_file.close()

    # Run solution to generate output
    os.system("./%s < Test%02d/%s.INP > Test%02d/%s.OUT" % (PROBLEM_NAME, test_id, PROBLEM_NAME, test_id, PROBLEM_NAME))