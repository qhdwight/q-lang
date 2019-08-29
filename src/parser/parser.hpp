#pragma once

#include <boost/program_options/variables_map.hpp>

#include <parser/node/master_node.hpp>

namespace po = boost::program_options;

namespace ql::parser {
    class Parser {
    private:
        using nodeFactory = std::function<std::shared_ptr<AbstractNode>(std::string&&, std::vector<std::string>&&, AbstractNode::ParentRef)>;

        std::map<std::string, nodeFactory> m_NamesToNodes;

        template<typename TNode>
        void registerNode(std::string_view nodeName) {
            m_NamesToNodes.emplace(nodeName, [](auto name, auto tokens, auto parent) {
                return std::make_shared<TNode>(std::move(name), std::move(tokens), parent);
            });
        }

        std::shared_ptr<AbstractNode> getNode
                (std::string const& nodeName, std::string&& blockWithInfo, std::vector<std::string>&& tokens, AbstractNode::ParentRef parent);

        void recurseNodes(std::string const& code, std::weak_ptr<AbstractNode> const& parent, int depth = 0);

    public:
        Parser();

        std::shared_ptr<MasterNode> parse(po::variables_map& options);

        std::shared_ptr<MasterNode> getNodes(std::string code);
    };
}
