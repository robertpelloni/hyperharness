from setuptools import setup, find_packages

setup(
    name="borg-sdk",
    version="0.1.0",
    packages=find_packages(),
    install_requires=[
        "requests>=2.31.0",
    ],
    author="OhMyOpenCode",
    description="Python SDK for Borg AI Orchestration System",
    keywords="ai, agents, orchestration, automation",
    python_requires=">=3.8",
)
